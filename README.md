# War Report Intelligence

## Descripción del Proyecto

**War Report** es un sistema automatizado de inteligencia de fuentes abiertas (OSINT) diseñado para monitorear y analizar de forma continua los principales conflictos armados globales. 

Utilizando el lenguaje **Go** y la **API de Google Gemini**, el sistema busca información en tiempo real enfocándose especialmente en los eventos de las **últimas 12 horas**. Su foco principal es rastrear:
- Armamento novedoso y nuevas tácticas en el campo de batalla.
- El uso de tecnologías y armas no convencionales por parte de actores irregulares, terroristas o insurgentes.
- Pruebas de armamento no declaradas u oficiales por parte de ejércitos regulares.
- Actualizaciones generales de los conflictos bélicos más relevantes.

El resultado de este análisis se inyecta dinámicamente en un archivo `index.html`, sirviendo como una fuente de noticias web estática que es auto-gestionada por el propio repositorio.

## Arquitectura y Estructura del Código

El proyecto sigue una estructura limpia y modularizada en Go, dividiendo responsabilidades en diferentes paquetes dentro del directorio `internal/`:

- **`main.go`**: Punto de entrada de la aplicación. Orquesta el flujo de: obtener el contexto previo, llamar al análisis de IA, y actualizar el sitio web.
- **`internal/gemini/client.go`**: Módulo principal de IA. Contiene la integración con el SDK de **Google GenAI**, las herramientas de búsqueda (Google Search Tool), las directrices del sistema y un mecanismo avanzado de **reintentos (Retry Logic)** con posibilidad de *fallback* a modelos secundarios para tolerar caídas o latencias en la API.
- **`internal/storage/`**: Gestiona la persistencia de los datos. Se encarga de leer el estado anterior en el archivo `index.html` y de reescribirlo con el nuevo informe HTML estructurado.
- **`internal/config/`**: Centraliza parámetros críticos como el ID del modelo principal, el modelo de fallback, delays para los reintentos y las instrucciones de sistema (*System Instructions*).
- **`internal/templates/`**: Maneja la estructura y las plantillas de visualización inyectables en el frontend.
- **`.github/workflows/monitor.yml`**: Pipeling de Integración Continua (CI/CD) que permite que el reporte viva de forma autónoma.

## Funcionalidades Clave

1. **Automatización Cron**: Se ejecuta automáticamente cada 12 horas gracias a GitHub Actions, sin intervención humana.
2. **Contexto Histórico Aware**: El sistema no analiza los eventos aislados. Antes de solicitar datos, recaba el "Reporte Anterior" desde el HTML y se lo proporciona a la IA, obligándola a buscar y resaltar **diferencias y avances** respecto al periodo pasado, evitando redundancias.
3. **Resiliencia (Fault Tolerance)**: Frente a posibles errores 503, 504 o interrupciones internas del servicio por parte de la API de Gemini, el sistema reintenta la conexión utilizando tiempos de espera programados y modelos de repuesto.
4. **GitOps Autónomo**: Una vez finalizado el ciclo, si hay nueva información confirmada, el sistema hace un commit automático y "pushea" el nuevo `index.html` al repositorio, actualizando inmediatamente cualquier sitio alojado vía GitHub Pages.

## Uso y Configuración Local

Para ejecutar o probar el código en un entorno de desarrollo local, se requiere:

1. **Go 1.22** o superior instalado.
2. Contar con una clave API válida de Google Gemini.
3. Crear un archivo `.env` en la raíz del proyecto o exportar la variable en la terminal:

```bash
export GEMINI_API_KEY="tu-clave-api"
```

### Ejecutar el reporte

```bash
go run main.go
```

En la consola se indicará si ha reconocido el contexto del archivo parcial. A continuación, el agente procesará el reporte y, si finaliza de manera exitosa, modificará localmente el `index.html` listando las referencias encontradas para mantener un buen nivel de SEO y corroboración.

## Configuración en GitHub (Despliegue)

Para asegurar el funcionamiento en la nube, es vital configurar en los ajustes del repositorio (Settings > Secrets and variables > Actions):
- El Secreto: `GEMINI_API_KEY` con un valor activo de API.

Además, el workflow de dependencias necesita permisos explicitados de escritura (`permissions: contents: write`) para publicar el commit que actualiza el portal en vivo.
