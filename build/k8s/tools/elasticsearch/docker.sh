docker pull elasticsearch:8.1.3
mkdir -p /opt/elasticsearch/config
mkdir -p /opt/elasticsearch/data
mkdir -p /opt/elasticsearch/plugins

echo "http.host: 0.0.0.0" >> ~/app/es/config/elasticsearch.yml

docker run --name elasticsearch \
 -e "discovery.type=single-node" \
 -e ES_JAVA_OPTS="-Xms84m -Xmx512m" \
 -v ~/app/es/config/elasticsearch.yml:/usr/share/elasticsearch/config/elasticsearch.yml \
 -v /data/es:/usr/share/elasticsearch/data \
 -v ~/app/es/plugins:/usr/share/elasticsearch/plugins \
 -d elasticsearch:7.17.3

 -m 1g \
 -p 9200:9200  -p 9300:9300 \

docker exec -it es01 /usr/share/elasticsearch/bin/elasticsearch-reset-password -u elastic
elasticsearch-setup-passwords interactive
docker cp es01:/usr/share/elasticsearch/config/certs/http_ca.crt .
curl --cacert http_ca.crt -u elastic https://localhost:9200
bin/elasticsearch-keystore show xpack.security.http.ssl.keystore.secure_password
bin/elasticsearch-keystore show xpack.security.transport.ssl.keystore.secure_password
docker exec -it es-node01 /usr/share/elasticsearch/bin/elasticsearch-create-enrollment-token -s kibana
