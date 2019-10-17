wget https://ftp.gnu.org/gnu/gcc/gcc-8.2.0/gcc-8.2.0.tar.gz
wget https://ftp.gnu.org/gnu/gmp/gmp-4.3.2.tar.gz
wget https://ftp.gnu.org/gnu/mpfr/mpfr-2.4.2.tar.gz
wget https://ftp.gnu.org/gnu/mpc/mpc-1.0.1.tar.gz
```bash
mkdir ~/local/gcc

tar xf gmp-4.3.2.tar.gz
cd gmp-4.3.2
sudo yum install -y m4
./configure --prefix=$HOME/local/gcc
make && make install
ls ~/local/gcc/lib/

cd ..
tar xf mpfr-2.4.2.tar.gz
cd mpfr-2.4.2
./configure --prefix=$HOME/local/gcc --with-gmp=$HOME/local/gcc
make && make install
ls ~/local/gcc/lib/

cd ..
tar xf mpc-1.0.1.tar.gz
cd mpc-1.0.1
./configure --prefix=$HOME/local/gcc --with-gmp=$HOME/local/gcc --with-mpfr=$HOME/local/gcc
make && make install
ls ~/local/gcc/lib/

cd ..
tar xf gcc-8.2.0.tar.gz
cd gcc-8.2.0

./configure --prefix=$HOME/local/gcc --with-gmp=$HOME/local/gcc --with-mpfr=$HOME/local/gcc --with-mpc=$HOME/local/gcc --disable-multilib
export LD_LIBRARY_PATH=$HOME/local/gcc/lib:$LD_LIBRARY_PATH
make && make install

export PATH=$HOME/local/gcc/bin:$PATH
```