package infrastructure

import (
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func InitLog(path string) (*zap.Logger, error) {
	viper.SetConfigFile(path)
	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	encoder := zapcore.EncoderConfig{
		MessageKey: viper.GetString("encoderConfig.messageKey"),
		CallerKey:  viper.GetString("encoderConfig.callerKey"),
		LevelKey:   viper.GetString("encoderConfig.levelKey"),
		TimeKey:    viper.GetString("encoderConfig.timeKey"),
	}

	_ = encoder.EncodeCaller.
		UnmarshalText([]byte(viper.GetString("encoderConfig.callerEncoder")))
	_ = encoder.EncodeLevel.
		UnmarshalText([]byte(viper.GetString("encoderConfig.levelEncoder")))
	_ = encoder.EncodeTime.
		UnmarshalText([]byte(viper.GetString("encoderConfig.timeEncoder")))

	zap.NewAtomicLevelAt(zapcore.DebugLevel)
	conf := zap.Config{
		Encoding:         viper.GetString("encoding"),
		OutputPaths:      viper.GetStringSlice("outputPaths"),
		ErrorOutputPaths: viper.GetStringSlice("errorOutputPaths"),
		EncoderConfig:    encoder,
	}
	if err := conf.Level.
		UnmarshalText([]byte(viper.GetString("level"))); err != nil {
		return nil, err
	}

	return conf.Build()
}
