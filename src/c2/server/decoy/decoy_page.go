package decoy

const DecoyHTML = `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>NexaCloud - Edge Delivery Network</title>
    <style>
        body {
            font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, Helvetica, Arial, sans-serif;
            background: #f8fafc;
            display: flex;
            justify-content: center;
            align-items: center;
            height: 100vh;
            margin: 0;
            padding: 20px;
            color: #1e293b;
        }
        .container {
            max-width: 700px;
            background: white;
            padding: 50px 60px;
            border-radius: 24px;
            box-shadow: 0 20px 60px rgba(0,0,0,0.08);
            text-align: center;
            border: 1px solid #e9eef3;
        }
        .badge {
            display: inline-block;
            background: #f1f5f9;
            padding: 6px 20px;
            border-radius: 40px;
            font-size: 13px;
            font-weight: 500;
            color: #475569;
            letter-spacing: 0.3px;
            margin-bottom: 20px;
        }
        h1 {
            font-size: 28px;
            font-weight: 600;
            margin: 0 0 12px 0;
            color: #0f172a;
        }
        .subtitle {
            font-size: 18px;
            color: #475569;
            margin: 0 0 24px 0;
            line-height: 1.6;
        }
        .status-box {
            background: #fef9e7;
            border: 1px solid #fde68a;
            border-radius: 12px;
            padding: 18px 24px;
            margin: 24px 0 20px 0;
            font-size: 15px;
            color: #92400e;
        }
        .ref {
            font-size: 13px;
            color: #94a3b8;
            margin: 24px 0 0 0;
            letter-spacing: 0.2px;
        }
        .footer {
            margin-top: 32px;
            padding-top: 20px;
            border-top: 1px solid #f1f5f9;
            font-size: 13px;
            color: #94a3b8;
        }
    </style>
</head>
<body>
    <div class="container">
        <div class="badge">EDGE DELIVERY NETWORK</div>
        <h1>Service Maintenance</h1>
        <p class="subtitle">
            We are currently performing scheduled upgrades to our content delivery infrastructure.
            All services will be restored shortly.
        </p>
        <div class="status-box">
            <strong>Status:</strong> Maintenance Window Active &bull; Expected completion in 15-30 minutes
        </div>
        <p class="ref">Reference: CE-2026-MAINT-009 &bull; Incident ID: NX-EDGE-2847</p>
        <div class="footer">
            &copy; 2026 NexaCloud &bull; <span style="color:#cbd5e1;">v3.2.1</span>
        </div>
    </div>
</body>
</html>`
