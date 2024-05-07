

#include <stdio.h>
#include <math.h>
#include <windows.h>
#include <tchar.h>

float f(float x, float y, float z) {
	float a = x * x + 9.0f / 4.0f * y * y + z * z - 1;
	return a * a * a - x * x * z * z * z - 9.0f / 80.0f * y * y * z * z * z;
}

float h(float x, float z) {
	for (float y = 1.0f; y >= 0.0f; y -= 0.001f)
		if (f(x, y, z) <= 0.0f)
			return y;
	return 0.0f;
}

int main() {
    // 定义控制台应用程序的入口点。
	HANDLE o = GetStdHandle(STD_OUTPUT_HANDLE);
	_TCHAR buffer[25][85] = { _T(' ') };
	for(int i = 0; i < 25; i++){
		_TCHAR* p = &buffer[i][0];
			*p++ = '\x1b';
			*p++ = '[';
			*p++ = '3';
			//*p++ = '1';
			*p++ = '1' + i%7;
			*p++ = 'm';
	}
	int times = 0;
	_TCHAR ramp[] = _T(".:-=+*#%@");
	for (float t = 0.0f;; t += 0.1f) {
		times++;
		if (times == 5){
			for(int i = 0; i < 25; i++){
			buffer[i][3] = (buffer[i][3] + 1);
			if(buffer[i][3] > '7'){
				buffer[i][3] = '1';
			}
			times = 0;
			}
		}
		int sy = 0;
		float s = sinf(t);
		float a = s * s * s * s * 0.2f;
		for (float z = 1.3f; z > -1.2f; z -= 0.1f) {
			_TCHAR* p = &buffer[sy++][5];
			float tz = z * (1.2f - a);
			for (float x = -1.5f; x < 1.5f; x += 0.05f) {
				float tx = x * (1.2f + a);
				float v = f(tx, 0.0f, tz);
				if (v <= 0.0f) {
					float y0 = h(tx, tz);
					float ny = 0.01f;
					float nx = h(tx + ny, tz) - y0;
					float nz = h(tx, tz + ny) - y0;
					float nd = 1.0f / sqrtf(nx * nx + ny * ny + nz * nz);
					float d = (nx + ny - nz) * nd * 0.5f + 0.5f;
					*p++ = ramp[(int)(d * 5.0f)];
				}
				else
					*p++ = ' ';
			}
		}
		for (sy = 0; sy < 25; sy++) {
			_TCHAR* p = &buffer[sy][80];
				*p++ = '\x1b';
				*p++ = '[';
				*p++ = '0';
				*p++ = 'm';
			COORD coord = { 0, sy };
			SetConsoleCursorPosition(o, coord);
			WriteConsole(o, buffer[sy], 85, NULL, 0);
		}
		Sleep(16);
	}
}
