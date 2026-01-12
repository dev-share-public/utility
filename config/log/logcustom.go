package logcustom

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"sync"
	"time"
	"utility/config/envi"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

type CallerInfo struct {
	Level    int    `json:"level"`
	File     string `json:"file"`
	Function string `json:"func"`
	Line     int    `json:"line"`
}

// getCallerInfos ดึงข้อมูล caller ตามระดับที่ระบุ (เช่น 2,3)
func getCallerInfos(levels ...int) []CallerInfo {
	var infos []CallerInfo
	for _, level := range levels {
		pc, file, line, ok := runtime.Caller(level)
		funcName := "unknown"
		if ok {
			fn := runtime.FuncForPC(pc)
			if fn != nil {
				funcName = fn.Name()
			}
		} else {
			file = "unknown"
			line = 0
		}
		infos = append(infos, CallerInfo{
			Level:    level,
			File:     file,
			Function: funcName,
			Line:     line,
		})
	}
	return infos
}

// NewZapErrorHandler คืนค่า fiber.ErrorHandler ที่เขียน log error ลงไฟล์
func NewZapErrorHandler(logDir string) fiber.ErrorHandler {
	now := time.Now().Format("2006-01-02")
	logFile := fmt.Sprintf("%s/system_error_%s.log", logDir, now)

	if err := os.MkdirAll(filepath.Dir(logFile), os.ModePerm); err != nil {
		panic(err)
	}

	writer := &lumberjack.Logger{
		Filename:   logFile,
		MaxSize:    10,
		MaxBackups: 7,
		MaxAge:     14,
		Compress:   false,
	}

	encoderCfg := zap.NewProductionEncoderConfig()
	encoderCfg.TimeKey = "time"
	encoderCfg.EncodeTime = zapcore.ISO8601TimeEncoder

	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderCfg),
		zapcore.AddSync(writer),
		zapcore.ErrorLevel,
	)

	logger := zap.New(core)

	return func(c *fiber.Ctx, err error) error {
		var body string
		if c.Request().Body() != nil {
			body = string(c.Body())
		}

		callers := getCallerInfos(2, 3)
		errType := reflect.TypeOf(err).String()

		logger.Error("System error occurred",
			zap.String("method", c.Method()),
			zap.String("path", c.OriginalURL()),
			zap.String("ip", c.IP()),
			zap.Any("headers", c.GetReqHeaders()),
			zap.String("query", c.Context().QueryArgs().String()),
			zap.String("body", body),
			zap.String("error", err.Error()),
			zap.String("error_type", errType),
			zap.Array("callers", zapcore.ArrayMarshalerFunc(func(enc zapcore.ArrayEncoder) error {
				for _, caller := range callers {
					enc.AppendObject(zapcore.ObjectMarshalerFunc(func(objEnc zapcore.ObjectEncoder) error {
						objEnc.AddInt("level", caller.Level)
						objEnc.AddString("file", caller.File)
						objEnc.AddString("func", caller.Function)
						objEnc.AddInt("line", caller.Line)
						return nil
					}))
				}
				return nil
			})),
		)

		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Internal Server Error",
		})
	}
}

type RequestResponseLogger struct {
	logDir     string
	currentDay string
	logger     *zap.Logger
	mu         sync.RWMutex
}

// NewRequestResponseLogger returns Fiber middleware
func NewRequestResponseLogger(logDir string) fiber.Handler {
	l := &RequestResponseLogger{
		logDir: logDir,
	}
	l.buildLoggerRequestResponse()         // สร้าง logger ครั้งแรก
	go l.startDailyRotateRequestResponse() // เริ่ม goroutine auto rotate

	return func(c *fiber.Ctx) error {
		start := time.Now()
		err := c.Next()

		// --- Check if logging is enabled via .env ---
		if envi.GetEnv("IS_DISPLAY_LOG_REQUEST_RESPONSE", "false") != "true" {
			return err // skip logging
		}

		method := c.Method()
		path := c.Path()
		ip := c.IP()
		status := c.Response().StatusCode()
		latency := time.Since(start)

		// collect headers
		reqHeaders := make(map[string]string)
		c.Request().Header.VisitAll(func(k, v []byte) {
			reqHeaders[string(k)] = string(v)
		})
		respHeaders := make(map[string]string)
		c.Response().Header.VisitAll(func(k, v []byte) {
			respHeaders[string(k)] = string(v)
		})

		reqBody := c.Body()
		respBody := c.Response().Body()

		l.mu.RLock()
		defer l.mu.RUnlock()
		l.logger.Info("HTTP Request/Response",
			zap.String("method", method),
			zap.String("path", path),
			zap.String("ip", ip),
			zap.Int("status", status),
			zap.Duration("latency", latency),
			zap.Any("request_headers", reqHeaders),
			zap.ByteString("request_body", reqBody),
			zap.Any("response_headers", respHeaders),
			zap.ByteString("response_body", respBody),
		)

		return err
	}
}

// buildLoggerRequestResponse สร้าง logger ใหม่
func (l *RequestResponseLogger) buildLoggerRequestResponse() {
	if _, err := os.Stat(l.logDir); os.IsNotExist(err) {
		_ = os.MkdirAll(l.logDir, 0755)
	}

	now := time.Now().Format("2006-01-02")
	l.currentDay = now
	logFile := filepath.Join(l.logDir, fmt.Sprintf("request-response-%s.log", now))

	writer := zapcore.AddSync(&lumberjack.Logger{
		Filename:   logFile,
		MaxSize:    10, // MB ต่อไฟล์
		MaxBackups: 5,
		MaxAge:     30,    // วัน
		Compress:   false, // ไม่บีบอัดไฟล์เก่า
	})

	encoderCfg := zap.NewProductionEncoderConfig()
	encoderCfg.TimeKey = "timestamp"
	encoderCfg.EncodeTime = zapcore.ISO8601TimeEncoder

	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderCfg),
		writer,
		zap.InfoLevel,
	)

	l.mu.Lock()
	l.logger = zap.New(core, zap.AddCaller())
	l.mu.Unlock()
}

// startDailyRotateRequestResponse goroutine สำหรับสร้าง logger ใหม่ทุกวัน
func (l *RequestResponseLogger) startDailyRotateRequestResponse() {
	for {
		time.Sleep(time.Minute)
		current := time.Now().Format("2006-01-02")
		if current != l.currentDay {
			l.buildLoggerRequestResponse()
		}
	}
}
