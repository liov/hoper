```bash
wget https://openresty.org/download/openresty-1.15.8.2.tar.gz
tar -xzvf openresty-VERSION.tar.gz
cd openresty-VERSION/

apt install libssl-dev libpq-dev libpcre3-dev libxml2-dev libxslt-dev libgd2-dev libgeoip-dev

./configure --with-select_module\
           --with-poll_module \
            --with-threads \
            --with-luajit \
            --without-http_redis2_module \
            --with-http_iconv_module \
            --with-http_postgres_module \
            --with-file-aio \
            --with-ipv6 \
            --with-http_v2_module \
            --with-http_realip_module \
            --with-http_addition_module \
            --with-http_xslt_module \
            --with-http_xslt_module=dynamic \
            --with-http_image_filter_module \
            --with-http_image_filter_module=dynamic \
            --with-http_geoip_module \
            --with-http_geoip_module=dynamic \
            --with-http_sub_module  \
            --with-http_dav_module  \
            --with-http_flv_module \
            --with-http_mp4_module  \
            --with-http_gunzip_module   \
            --with-http_gzip_static_module  \
            --with-http_auth_request_module \
            --with-http_random_index_module  \
            --with-http_secure_link_module  \
            --with-http_degradation_module   \
            --with-http_slice_module       \
            --with-http_stub_status_module  \
            --with-mail  \
            --with-mail=dynamic   \
            --with-mail_ssl_module   \
            --with-stream      \
            --with-stream=dynamic   \
            --with-stream_ssl_module    \
            --with-stream_realip_module    \
            --with-stream_geoip_module    \
            --with-stream_geoip_module=dynamic \
            --with-stream_ssl_preread_module

make
sudo make install DESTDIR=/whereto make install
./nginx -c nginx.conf
./nginx -c nginx.conf -s reload
```
