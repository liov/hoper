-- 源码
wget https://nodejs.org/dist/v13.3.0/node-v13.3.0.tar.gz
tar -xzvf

apt install g++
sudo apt-get install python3-distutils
./configure --prefix=/usr/local/node
make && make install
-- 二进制
wget https://nodejs.org/dist/v12.3.1/node-v12.3.1-linux-x64.tar.xz

tar -Jxvf

wget https://nodejs.org/dist/v16.0.0/node-v16.0.0-linux-x64.tar.xz
export PATH=/home/node-v16.0.0-linux-x64/bin:$PATH

# nvm

nvm install 16.4.0 64
nvm use 16.4.0
