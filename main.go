package main

import (
	"context"
	"fmt"
	"log"

	"war-report/internal/gemini"
	"war-report/internal/storage"
)

func main() {
	ctx := context.Background()

	previousReport := storage.ExtractPreviousReport()
	if previousReport != "" {
		fmt.Println(">>> Reporte anterior encontrado. Se usará como contexto para comparación.")
	} else {
		fmt.Println(">>> No se encontró reporte anterior. Generando primer reporte.")
	}

	report, err := gemini.GetAnalysis(ctx, previousReport)
	if err != nil {
		log.Fatalf("ERROR FATAL: %v", err)
	}

	err = storage.UpdateHTML(report)
	if err != nil {
		log.Fatalf("ERROR FATAL actualizando HTML: %v", err)
	}

	fmt.Println("SUCCESS: HTML actualizado correctamente.")
}
