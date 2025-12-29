#!/bin/bash
# 启动 Elasticsearch 和 Kibana（用于日志功能）

echo "启动 Elasticsearch..."
docker run -d \
  --name talent-elasticsearch \
  -p 9200:9200 \
  -e "discovery.type=single-node" \
  -e "xpack.security.enabled=false" \
  -e "ES_JAVA_OPTS=-Xms512m -Xmx512m" \
  docker.elastic.co/elasticsearch/elasticsearch:8.11.0

echo "等待 Elasticsearch 启动..."
sleep 30

# 检查 ES 是否启动成功
if curl -s http://localhost:9200/_cluster/health | grep -q 'green\|yellow'; then
    echo "✓ Elasticsearch 启动成功"
else
    echo "✗ Elasticsearch 启动失败，请检查日志"
    exit 1
fi

echo ""
echo "启动 Kibana（可选，用于日志可视化）..."
docker run -d \
  --name talent-kibana \
  -p 5601:5601 \
  -e "ELASTICSEARCH_HOSTS=http://host.docker.internal:9200" \
  docker.elastic.co/kibana/kibana:8.11.0

echo ""
echo "Elasticsearch: http://localhost:9200"
echo "Kibana: http://localhost:5601 (启动需要1-2分钟)"
