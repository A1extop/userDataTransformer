package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"userDataTransformer/internal/config"
)

// настройка логгера
func SetupLogger(cfg *config.Config) (*zap.Logger, error) {
	zapCfg := zap.NewProductionConfig()

	level, err := zap.ParseAtomicLevel(cfg.Log.Level)
	if err != nil {
		return nil, err
	}

	zapCfg.Level = level
	zapCfg.Encoding = cfg.Log.Encoding

	var outputs []string
	if cfg.Log.EnableConsole {
		outputs = append(outputs, "stdout")
	}
	if cfg.Log.EnableFile {
		outputs = append(outputs, cfg.Log.FilePath)
	}

	if len(outputs) == 0 {
		outputs = []string{"stdout"}
	}

	zapCfg.OutputPaths = outputs
	zapCfg.ErrorOutputPaths = outputs

	zapCfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	zapCfg.DisableStacktrace = true

	return zapCfg.Build()
}
