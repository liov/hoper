vscode调试flask真坑
launch.json
```json
	"version": "0.2.0",
	"configurations": [
		{
			"name": "Python: Flask",
			"type": "python",
			"request": "launch",
			"module": "flask",
			"env": {
				"FLASK_APP": "${file}.py",
				"FLASK_DEBUG": "1"
			},
			"args": [
				"run",
				"--no-debugger",
				"--no-reload"
			],
			"jinja": true,
			"justMyCode": true
		}
	]
}
```