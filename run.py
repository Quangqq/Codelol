from flask import Flask, jsonify
import subprocess

app = Flask(__name__)

@app.route('/run-main', methods=['GET'])
def run_main():
    try:
        # Chạy main.py
        result = subprocess.run(['python', 'main.py'], capture_output=True, text=True)
        return jsonify({
            'status': 'success',
            'message': result.stdout
        })
    except Exception as e:
        return jsonify({
            'status': 'error',
            'message': str(e)
        })

if __name__ == '__main__':
    app.run(host='0.0.0.0', port=4000)  # Mở cổng 5000
