package gemini

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"war-report/internal/config"

	"github.com/joho/godotenv"
	"google.golang.org/genai"
)

type attemptInfo struct {
	model string
	label string
}

func GetAnalysis(ctx context.Context, previousReport string) (string, error) {
	apiKey := os.Getenv("GEMINI_API_KEY")
	if apiKey == "" {
		godotenv.Load()
		apiKey = os.Getenv("GEMINI_API_KEY")
	}

	if apiKey == "" {
		return "", fmt.Errorf("falta la API key. Define GEMINI_API_KEY o GOOGLE_API_KEY en las variables de entorno")
	}

	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		Backend: genai.BackendGeminiAPI,
		APIKey:  apiKey,
	})
	if err != nil {
		return "", fmt.Errorf("error inicializando cliente: %w", err)
	}

	nowStr := time.Now().UTC().Format("2006-01-02 15:04 UTC")
	previousContext := ""

	if previousReport != "" {
		previousContext = fmt.Sprintf(`
    --- CONTEXTO: REPORTE ANTERIOR ---
    A continuación se encuentra el reporte generado previamente. Úsalo como referencia para:
    1. Hacer una comparación de avance entre la nueva información que vas a encontrar en las últimas 12 horas y lo ya reportado anteriormente.
    2. Identificar novedades y desarrollos sin repetir la misma información, enfatizando lo que ha cambiado.
    
    REPORTE ANTERIOR:
    %s
    
    IMPORTANTE: Integra la comparación en el desarrollo de los puntos, resaltando explícitamente qué hay de nuevo respecto al reporte anterior.
    --- FIN DEL CONTEXTO ANTERIOR ---
    `, previousReport)
	}

	query := fmt.Sprintf(`
    CURRENT TIME: %s%s
    
    **TASK: Ejecuta el Informe OSINT de Conflictos y Novedades Armamentísticas.**
    
    Sigue estrictamente estos pasos de investigación usando Google Search (priorizando información de las ÚLTIMAS 12 HORAS):

    PASO 1: Búsqueda sobre Armamento Novedoso y Tácticas en Batalla
    - Busca noticias recientes sobre nuevas armas, prototipos o nuevas tácticas usadas en frentes de batalla (menos de 6 meses de existencia).
    
    PASO 2: Actores Irregulares y Terrorismo
    - Busca reportes recientes de grupos terroristas, insurgentes o al margen de la ley utilizando tecnologías o armas nuevas/novedosas en el mundo.
    
    PASO 3: Pruebas No Oficiales de Ejércitos
    - Investiga si hay indicios o denuncias sobre ejércitos regulares probando armas experimentales no declaradas oficialmente en el campo de batalla.
    
    PASO 4: Monitoreo de Principales Conflictos
    - Actualizaciones fundamentales y últimos movimientos bélicos globales.
    
    OUTPUT FORMAT (HTML puro):
    <div class="report-section">
      <h3>Estado de los Principales Conflictos (Últimas 12h)</h3>
      <p>[Actualizaciones más importantes, con comparativa de avance si existía reporte anterior]</p>
      
      <h3>Innovación en el Campo de Batalla (Armamento y Tácticas)</h3>
      <ul>
        <li><strong>[Arma/Táctica]:</strong> [Descripción de la novedad y dónde se está empleando]</li>
      </ul>
      
      <h3>Actores al Margen de la Ley</h3>
      <p>[Uso de armas novedosas por grupos terroristas y actores no estatales]</p>
      
      <h3>Laboratorios de Guerra (Pruebas No Oficiales)</h3>
      <p>[Información sobre ejércitos testeando armamento no revelado]</p>
      
      <h3>Fuentes Consultadas</h3>
      <ul>
        <!-- IMPORTANTE: Rellena esta lista ÚNICAMENTE con enlaces REALES obtenidos en tu búsqueda. PROHIBIDO INVENTAR URLs. SI NO ESTÁS SEGURO DE UN ENLACE, NO LO INCLUYAS. -->
        <li><a href="URL_REAL_OBTENIDA_EN_BUSQUEDA">Título de la Fuente Verificada</a></li>
      </ul>
    </div>
    `, nowStr, previousContext)

	tools := []*genai.Tool{
		{
			GoogleSearch: &genai.GoogleSearch{},
		},
	}

	genConfig := &genai.GenerateContentConfig{
		Tools: tools,
		SystemInstruction: &genai.Content{
			Parts: []*genai.Part{
				{Text: config.SysInstruct},
			},
		},
		Temperature: genai.Ptr[float32](0.3),
	}

	contents := []*genai.Content{
		{
			Role: "user",
			Parts: []*genai.Part{
				{Text: query},
			},
		},
	}

	attempts := []attemptInfo{
		{config.ModelID, "principal"},
		{config.ModelID, "principal"},
		{config.ModelIDFallback, "fallback"},
		{config.ModelIDFallback, "fallback"},
	}

	fmt.Println(">>>> Iniciando contacto con Gemini (Estratega Macro)...")

	var lastErr error
	for attemptIdx, attempt := range attempts {
		callCtx, cancel := context.WithTimeout(ctx, config.CallTimeout)
		result, err := client.Models.GenerateContent(callCtx, attempt.model, contents, genConfig)
		cancel()

		if err == nil {
			fmt.Println("\n>>> Análisis completado.")
			return result.Text(), nil
		}

		lastErr = err
		infoErr := err.Error()
		isTimeout := false
		isInternalError := false
		// check if infoErr contains 504
		if strings.Contains(infoErr, "504") {
			isTimeout = true
		}
		// check if infoErr contains 500
		if strings.Contains(infoErr, "500") {
			isInternalError = true
		}
		reason := "Error 503 o de servicio"
		if isTimeout {
			reason = fmt.Sprintf("Timeout (%ds)", int(config.CallTimeout.Seconds()))
		}
		if isInternalError {
			reason = "Error interno"
		}

		if attemptIdx < len(attempts)-1 {
			nextModelLabel := attempts[attemptIdx+1].label
			fmt.Printf("\n[ADVERTENCIA] %s en modelo %s (Intento %d/%d)\n", reason, attempt.label, attemptIdx+1, len(attempts))
			fmt.Printf(">>> Esperando %ds antes de reintentar con modelo %s...\n", int(config.RetryDelay.Seconds()), nextModelLabel)
			time.Sleep(config.RetryDelay)
		} else {
			fmt.Printf("\n[ERROR] %s persistente en todos los modelos (%d intentos).\n", strings.ToLower(reason), len(attempts))
			fmt.Println(">>> Terminando ejecución para evitar costos o bucles infinitos.")
		}
	}

	return "", fmt.Errorf("API de Gemini no disponible después de %d intentos, ultimo error: %w", len(attempts), lastErr)
}
