package templates

const BaseHTML = `<!DOCTYPE html>
<html lang="es">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <meta name="theme-color" content="#0f172a">
    <title>WAR Monitor - Inteligencia Geopolítica</title>
    <style>
        :root {
            --bg-body: #0f172a;       /* Slate 900 */
            --bg-card: #1e293b;       /* Slate 800 */
            --text-main: #f1f5f9;     /* Slate 100 */
            --text-muted: #94a3b8;    /* Slate 400 */
            --border: #334155;        /* Slate 700 */
            --accent: #38bdf8;        /* Sky 400 */
            --accent-glow: rgba(56, 189, 248, 0.15);
            --success: #34d399;       /* Emerald 400 */
            --font-main: 'Inter', 'Segoe UI', system-ui, -apple-system, sans-serif;
        }

        body {
            font-family: var(--font-main);
            background-color: var(--bg-body);
            color: var(--text-main);
            margin: 0;
            padding: 20px;
            line-height: 1.6;
            -webkit-font-smoothing: antialiased;
        }

        .container {
            max-width: 850px;
            margin: 0 auto;
        }

        .header {
            text-align: center;
            margin-bottom: 40px;
            padding-bottom: 25px;
            border-bottom: 1px solid var(--border);
            animation: fadeIn 0.8s ease-out;
        }

        .header h1 {
            margin: 0 0 10px 0;
            font-size: 2rem;
            font-weight: 700;
            letter-spacing: -0.02em;
            background: linear-gradient(to right, #fff, #cbd5e1);
            -webkit-background-clip: text;
            -webkit-text-fill-color: transparent;
        }

        .header p {
            margin: 0;
            color: var(--text-muted);
            font-size: 0.95rem;
        }

        /* Tarjetas de Reporte */
        .report-card {
            background-color: var(--bg-card);
            border: 1px solid var(--border);
            border-radius: 16px;
            padding: 24px;
            margin-bottom: 24px;
            box-shadow: 0 4px 6px -1px rgba(0, 0, 0, 0.2);
            position: relative;
            overflow: hidden;
            transition: transform 0.2s, box-shadow 0.2s;
        }

        .report-card:hover {
            transform: translateY(-2px);
            box-shadow: 0 10px 15px -3px rgba(0, 0, 0, 0.3);
            border-color: var(--accent);
        }

        /* Borde lateral de acento */
        .report-card::before {
            content: "";
            position: absolute;
            left: 0; top: 0; bottom: 0;
            width: 4px;
            background: var(--accent);
            opacity: 0.8;
        }

        .timestamp {
            display: inline-block;
            font-size: 0.75rem;
            font-weight: 700;
            text-transform: uppercase;
            letter-spacing: 0.05em;
            color: var(--accent);
            background: var(--accent-glow);
            padding: 4px 12px;
            border-radius: 99px;
            margin-bottom: 20px;
            border: 1px solid rgba(56, 189, 248, 0.2);
        }

        /* Estilos del contenido generado */
        .content h3 {
            color: var(--text-main);
            font-size: 1.1rem;
            margin-top: 24px;
            margin-bottom: 12px;
            padding-bottom: 8px;
            border-bottom: 1px solid var(--border);
            display: flex;
            align-items: center;
        }
        
        .content h3::before {
            content: "▹";
            margin-right: 8px;
            color: var(--accent);
        }

        .content p {
            color: var(--text-muted);
            margin-bottom: 16px;
        }

        .content a {
            color: var(--accent);
            text-decoration: none;
            transition: color 0.2s ease;
        }
        
        .content a:hover {
            color: #bae6fd; 
            text-decoration: underline;
        }

        .content strong {
            color: var(--success); /* Resalta datos clave en verde */
            font-weight: 600;
        }
        
        /* Diagnóstico final destacado */
        .content p:last-child strong {
            color: #fbbf24; /* Amber para la conclusion */
        }

        @keyframes fadeIn {
            from { opacity: 0; transform: translateY(10px); }
            to { opacity: 1; transform: translateY(0); }
        }

        /* Mobile */
        @media (max-width: 600px) {
            body { padding: 12px; }
            .header h1 { font-size: 1.5rem; }
            .report-card { padding: 20px; border-radius: 12px; }
        }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>WAR Monitor {current_year}</h1>
            <p>Análisis de Riesgos Geopolíticos &amp; Tecnológicos Estratégicos</p>
        </div>
        
        <div id="archive">
            <!-- MR_REPORTS_START -->
            <!-- MR_REPORTS_END -->
        </div>
    </div>
    <script>
        function updateTimestamps() {
            const elements = document.querySelectorAll('.timestamp');
            const now = new Date();

            elements.forEach(el => {
                if (!el.dataset.originalTime) {
                    const text = el.textContent;
                    const match = text.match(/REPORTE GENERADO:\s*([^<]+)/);
                    if (match && match[1]) {
                        // Sometimes the match could include previously appended "(hace...)", 
                        // but since we only do this once per element (guarded by dataset.originalTime),
                        // and the server doesn't save the JS-modified DOM, the base text is clean.
                        let timeStr = match[1].trim();
                        // Strip out any existing "(hace...)" just in case
                        timeStr = timeStr.replace(/\s*\(hace.*\)$/, '');
                        
                        el.dataset.originalTime = timeStr;
                        // Base text should also not have any "(hace...)"
                        el.dataset.baseText = text.replace(/\s*\(hace.*\)$/, '').trim();
                    }
                }

                if (el.dataset.originalTime) {
                    const reportDate = new Date(el.dataset.originalTime);
                    if (!isNaN(reportDate)) {
                        const diffMs = now - reportDate;
                        if (diffMs < 0) return; // if report is somehow in the future

                        let diffMins = Math.floor(diffMs / 60000);
                        let diffHours = Math.floor(diffMins / 60);
                        const diffDays = Math.floor(diffHours / 24);

                        diffHours = diffHours % 24;
                        diffMins = diffMins % 60;

                        let timeParts = [];
                        if (diffDays > 0) {
                            timeParts.push(diffDays + ' día' + (diffDays !== 1 ? 's' : ''));
                        }
                        if (diffHours > 0) {
                            timeParts.push(diffHours + ' hora' + (diffHours !== 1 ? 's' : ''));
                        }
                        if (diffMins > 0) {
                            timeParts.push(diffMins + ' minuto' + (diffMins !== 1 ? 's' : ''));
                        }

                        let relativeTime = '';
                        if (timeParts.length === 0) {
                            relativeTime = 'hace unos segundos';
                        } else if (timeParts.length === 1) {
                            relativeTime = 'hace ' + timeParts[0];
                        } else if (timeParts.length === 2) {
                            relativeTime = 'hace ' + timeParts[0] + ' y ' + timeParts[1];
                        } else if (timeParts.length === 3) {
                            relativeTime = 'hace ' + timeParts[0] + ', ' + timeParts[1] + ' y ' + timeParts[2];
                        }

                        el.textContent = el.dataset.baseText + ' (' + relativeTime + ')';
                    }
                }
            });
        }

        document.addEventListener('DOMContentLoaded', () => {
            updateTimestamps();
            setInterval(updateTimestamps, 60000);
        });
    </script>
</body>
</html>`
