from setuptools import setup,Extension

MOD = 'CMOD'
setup(name=MOD, ext_modules=[Extension(MOD, sources=['callc.cpp'])])