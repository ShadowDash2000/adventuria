package main

import (
	"errors"
	"log/slog"
	"os"
	"path/filepath"
	"text/template"
)

const (
	effectOutputDir    = "././internal/adventuria_new/effects/custom/"
	effectTemplatePath = "././internal/adventuria_new/effects/effect.tmpl"
)

type createEffectData struct {
	Name        string
	RawName     string
	FirstLetter string
}

func createEffect() error {
	if len(os.Args) < 4 {
		slog.Error("Usage: <command> <type> <effect_name>")
		os.Exit(1)
	}

	data := createEffectData{
		Name:        snakeToCamelCase(os.Args[3]),
		RawName:     os.Args[3],
		FirstLetter: os.Args[3][:1],
	}

	fileName := data.RawName + ".go"
	outputDir := effectOutputDir + data.RawName

	err := os.MkdirAll(outputDir, os.ModePerm)
	if err != nil {
		return err
	}

	filePath := filepath.Join(outputDir, fileName)

	_, err = os.Stat(filePath)
	if err == nil {
		return errors.New("file already exists")
	}

	tmplBytes, err := os.ReadFile(effectTemplatePath)
	if err != nil {
		return err
	}

	tmpl, err := template.New("effect").Parse(string(tmplBytes))
	if err != nil {
		return err
	}

	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	err = tmpl.Execute(file, data)
	if err != nil {
		return err
	}

	slog.Info("Effect created successfully: " + filePath)

	return nil
}
