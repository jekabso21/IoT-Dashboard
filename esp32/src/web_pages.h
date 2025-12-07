#ifndef WEB_PAGES_H
#define WEB_PAGES_H

const char WIFI_PAGE[] PROGMEM = R"rawliteral(
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>WiFi Setup</title>
    <style>
        * {
            margin: 0;
            padding: 0;
            box-sizing: border-box;
        }
        body {
            font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Oxygen, Ubuntu, Cantarell, sans-serif;
            background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
            min-height: 100vh;
            display: flex;
            justify-content: center;
            align-items: center;
            padding: 20px;
        }
        .container {
            background: white;
            border-radius: 16px;
            box-shadow: 0 20px 60px rgba(0,0,0,0.3);
            max-width: 480px;
            width: 100%;
            padding: 40px 30px;
        }
        h1 {
            color: #333;
            font-size: 28px;
            margin-bottom: 10px;
            text-align: center;
        }
        .subtitle {
            color: #666;
            font-size: 14px;
            text-align: center;
            margin-bottom: 30px;
        }
        .form-group {
            margin-bottom: 20px;
        }
        label {
            display: block;
            color: #555;
            font-weight: 600;
            margin-bottom: 8px;
            font-size: 14px;
        }
        select, input {
            width: 100%;
            padding: 12px 15px;
            border: 2px solid #e0e0e0;
            border-radius: 8px;
            font-size: 16px;
            transition: border-color 0.3s;
        }
        select:focus, input:focus {
            outline: none;
            border-color: #667eea;
        }
        .btn {
            width: 100%;
            padding: 14px;
            background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
            color: white;
            border: none;
            border-radius: 8px;
            font-size: 16px;
            font-weight: 600;
            cursor: pointer;
            transition: transform 0.2s, box-shadow 0.2s;
            margin-top: 10px;
        }
        .btn:hover {
            transform: translateY(-2px);
            box-shadow: 0 5px 15px rgba(102, 126, 234, 0.4);
        }
        .btn:disabled {
            opacity: 0.6;
            cursor: not-allowed;
            transform: none;
        }
        .status {
            margin-top: 20px;
            padding: 15px;
            border-radius: 8px;
            font-size: 14px;
            display: none;
        }
        .status.success {
            background: #d4edda;
            color: #155724;
            border: 1px solid #c3e6cb;
        }
        .status.error {
            background: #f8d7da;
            color: #721c24;
            border: 1px solid #f5c6cb;
        }
        .status.info {
            background: #d1ecf1;
            color: #0c5460;
            border: 1px solid #bee5eb;
        }
        .loader {
            border: 3px solid #f3f3f3;
            border-top: 3px solid #667eea;
            border-radius: 50%;
            width: 20px;
            height: 20px;
            animation: spin 1s linear infinite;
            display: inline-block;
            margin-right: 10px;
        }
        @keyframes spin {
            0% { transform: rotate(0deg); }
            100% { transform: rotate(360deg); }
        }
        .network-count {
            color: #888;
            font-size: 12px;
            margin-top: 5px;
        }
    </style>
</head>
<body>
    <div class="container">
        <h1>IoT Dashboard Setup</h1>
        <p class="subtitle">Configure your device's WiFi connection</p>

        <form id="wifiForm">
            <div class="form-group">
                <label for="ssid">WiFi Network</label>
                <select id="ssid" required>
                    <option value="">Scanning networks...</option>
                </select>
                <div class="network-count" id="networkCount"></div>
            </div>

            <div class="form-group">
                <label for="password">Password</label>
                <input type="password" id="password" placeholder="Enter WiFi password">
            </div>

            <button type="submit" class="btn" id="connectBtn">Connect to WiFi</button>
            <button type="button" class="btn" id="refreshBtn" onclick="scanNetworks()" style="background: #6c757d;">Refresh Networks</button>
        </form>

        <div id="status" class="status"></div>
    </div>

    <script>
        function showStatus(message, type) {
            const status = document.getElementById('status');
            status.textContent = message;
            status.className = 'status ' + type;
            status.style.display = 'block';
        }

        function showStatusWithLoader(message, type) {
            const status = document.getElementById('status');
            status.innerHTML = '<div class="loader"></div>' + message;
            status.className = 'status ' + type;
            status.style.display = 'block';
        }

        async function scanNetworks() {
            const ssidSelect = document.getElementById('ssid');
            const refreshBtn = document.getElementById('refreshBtn');
            const networkCount = document.getElementById('networkCount');

            ssidSelect.disabled = true;
            refreshBtn.disabled = true;
            ssidSelect.innerHTML = '<option>Scanning...</option>';

            try {
                const response = await fetch('/scan');
                const data = await response.json();
                const networks = data.networks || [];

                ssidSelect.innerHTML = '';
                if (networks.length === 0) {
                    ssidSelect.innerHTML = '<option value="">No networks found</option>';
                } else {
                    networks.forEach(network => {
                        const option = document.createElement('option');
                        option.value = network.ssid;
                        const lock = network.encrypted ? ' [Secured]' : ' [Open]';
                        option.textContent = network.ssid + ' (' + network.rssi + ' dBm)' + lock;
                        ssidSelect.appendChild(option);
                    });
                    networkCount.textContent = 'Found ' + networks.length + ' networks';
                }
            } catch (error) {
                ssidSelect.innerHTML = '<option value="">Scan failed</option>';
                showStatus('Failed to scan: ' + error.message, 'error');
            } finally {
                ssidSelect.disabled = false;
                refreshBtn.disabled = false;
            }
        }

        document.getElementById('wifiForm').addEventListener('submit', async (e) => {
            e.preventDefault();

            const ssid = document.getElementById('ssid').value;
            const password = document.getElementById('password').value;
            const connectBtn = document.getElementById('connectBtn');

            if (!ssid) {
                showStatus('Please select a WiFi network', 'error');
                return;
            }

            connectBtn.disabled = true;
            showStatusWithLoader('Connecting to WiFi...', 'info');

            try {
                const response = await fetch('/connect', {
                    method: 'POST',
                    headers: { 'Content-Type': 'application/json' },
                    body: JSON.stringify({ ssid, password })
                });

                const data = await response.json();

                if (data.success) {
                    showStatus('Connected! Redirecting to pairing...', 'success');
                    setTimeout(() => { window.location.href = '/pairing'; }, 2000);
                } else {
                    showStatus('Failed: ' + (data.message || 'Unknown error'), 'error');
                    connectBtn.disabled = false;
                }
            } catch (error) {
                showStatus('Error: ' + error.message, 'error');
                connectBtn.disabled = false;
            }
        });

        scanNetworks();
    </script>
</body>
</html>
)rawliteral";

const char PAIRING_PAGE[] PROGMEM = R"rawliteral(
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Device Pairing</title>
    <style>
        * {
            margin: 0;
            padding: 0;
            box-sizing: border-box;
        }
        body {
            font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Oxygen, Ubuntu, Cantarell, sans-serif;
            background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
            min-height: 100vh;
            display: flex;
            justify-content: center;
            align-items: center;
            padding: 20px;
        }
        .container {
            background: white;
            border-radius: 16px;
            box-shadow: 0 20px 60px rgba(0,0,0,0.3);
            max-width: 480px;
            width: 100%;
            padding: 40px 30px;
            text-align: center;
        }
        h1 {
            color: #333;
            font-size: 28px;
            margin-bottom: 10px;
        }
        .subtitle {
            color: #666;
            font-size: 14px;
            margin-bottom: 30px;
        }
        .instructions {
            background: #f8f9fa;
            border-left: 4px solid #667eea;
            padding: 15px;
            margin-bottom: 30px;
            text-align: left;
            border-radius: 4px;
        }
        .instructions h3 {
            color: #333;
            font-size: 16px;
            margin-bottom: 10px;
        }
        .instructions ol {
            color: #666;
            font-size: 14px;
            padding-left: 20px;
        }
        .code-inputs {
            display: flex;
            justify-content: center;
            gap: 10px;
            margin-bottom: 20px;
        }
        .code-digit {
            width: 50px;
            height: 60px;
            font-size: 32px;
            font-weight: bold;
            text-align: center;
            border: 2px solid #e0e0e0;
            border-radius: 8px;
        }
        .code-digit:focus {
            outline: none;
            border-color: #667eea;
        }
        .btn {
            width: 100%;
            padding: 14px;
            background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
            color: white;
            border: none;
            border-radius: 8px;
            font-size: 16px;
            font-weight: 600;
            cursor: pointer;
            margin-top: 10px;
        }
        .btn:disabled {
            opacity: 0.6;
            cursor: not-allowed;
        }
        .btn-secondary {
            background: #6c757d;
        }
        .status {
            margin-top: 20px;
            padding: 15px;
            border-radius: 8px;
            font-size: 14px;
            display: none;
        }
        .status.success { background: #d4edda; color: #155724; }
        .status.error { background: #f8d7da; color: #721c24; }
        .status.info { background: #d1ecf1; color: #0c5460; }
        .loader {
            border: 3px solid #f3f3f3;
            border-top: 3px solid #667eea;
            border-radius: 50%;
            width: 20px;
            height: 20px;
            animation: spin 1s linear infinite;
            display: inline-block;
            margin-right: 10px;
        }
        @keyframes spin {
            0% { transform: rotate(0deg); }
            100% { transform: rotate(360deg); }
        }
    </style>
</head>
<body>
    <div class="container">
        <h1>Device Pairing</h1>
        <p class="subtitle">Connect your device to IoT Dashboard</p>

        <div class="instructions">
            <h3>How to pair:</h3>
            <ol>
                <li>Open IoT Dashboard web app</li>
                <li>Go to Devices and click "Add Device"</li>
                <li>Copy the 6-digit pairing code</li>
                <li>Enter the code below</li>
            </ol>
        </div>

        <form id="pairingForm">
            <div class="code-inputs">
                <input type="text" class="code-digit" maxlength="1" pattern="[0-9]" required>
                <input type="text" class="code-digit" maxlength="1" pattern="[0-9]" required>
                <input type="text" class="code-digit" maxlength="1" pattern="[0-9]" required>
                <input type="text" class="code-digit" maxlength="1" pattern="[0-9]" required>
                <input type="text" class="code-digit" maxlength="1" pattern="[0-9]" required>
                <input type="text" class="code-digit" maxlength="1" pattern="[0-9]" required>
            </div>

            <button type="submit" class="btn" id="pairBtn">Pair Device</button>
            <button type="button" class="btn btn-secondary" onclick="window.location.href='/'">Back to WiFi</button>
        </form>

        <div id="status" class="status"></div>
    </div>

    <script>
        const codeInputs = document.querySelectorAll('.code-digit');

        codeInputs.forEach((input, index) => {
            input.addEventListener('input', (e) => {
                if (e.target.value.length === 1 && index < codeInputs.length - 1) {
                    codeInputs[index + 1].focus();
                }
            });

            input.addEventListener('keydown', (e) => {
                if (e.key === 'Backspace' && !e.target.value && index > 0) {
                    codeInputs[index - 1].focus();
                }
            });

            input.addEventListener('keypress', (e) => {
                if (!/[0-9]/.test(e.key)) e.preventDefault();
            });
        });

        codeInputs[0].focus();

        function showStatus(message, type) {
            const status = document.getElementById('status');
            status.textContent = message;
            status.className = 'status ' + type;
            status.style.display = 'block';
        }

        document.getElementById('pairingForm').addEventListener('submit', async (e) => {
            e.preventDefault();

            const code = Array.from(codeInputs).map(i => i.value).join('');
            if (code.length !== 6) {
                showStatus('Please enter all 6 digits', 'error');
                return;
            }

            const pairBtn = document.getElementById('pairBtn');
            pairBtn.disabled = true;
            document.getElementById('status').innerHTML = '<div class="loader"></div>Pairing...';
            document.getElementById('status').className = 'status info';
            document.getElementById('status').style.display = 'block';

            try {
                const response = await fetch('/pair', {
                    method: 'POST',
                    headers: { 'Content-Type': 'application/json' },
                    body: JSON.stringify({ code })
                });

                const data = await response.json();

                if (data.success) {
                    showStatus('Device paired successfully! Starting normal operation...', 'success');
                } else {
                    showStatus('Failed: ' + (data.message || 'Invalid code'), 'error');
                    pairBtn.disabled = false;
                    codeInputs.forEach(i => { i.value = ''; });
                    codeInputs[0].focus();
                }
            } catch (error) {
                showStatus('Error: ' + error.message, 'error');
                pairBtn.disabled = false;
            }
        });
    </script>
</body>
</html>
)rawliteral";

#endif
