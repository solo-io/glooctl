package docker

const (
	envoyYAML = `#envoy.yaml
node:
  cluster: ingress
  id: ingress-1
static_resources:
  clusters:
  - name: xds_cluster
    connect_timeout: 5.000s
    hosts:
    - socket_address:
        address: control-plane
        port_value: 8081
    http2_protocol_options: {}
    type: STRICT_DNS
dynamic_resources:
  ads_config:
    api_type: GRPC
    cluster_names:
    - xds_cluster
  cds_config:
    ads: {}
  lds_config:
    ads: {}
admin:
  access_log_path: /dev/stdout
  address:
    socket_address:
      address: 0.0.0.0
      port_value: 19000`

	dockerComposeYAML = `# docker compose for gloo
version: '3'
services:
  ingress:
    image: soloio/envoy:v0.1.6-127
    entrypoint: ["envoy"]
    command: ["-c", "/config/envoy.yaml", "--v2-config-only"]
    volumes:
    - ./envoy.yaml:/config/envoy.yaml:ro
    ports:
    - "8080:8080"
    - "19000:19000"

  control-plane:
    image: soloio/control-plane:0.2.0
    entrypoint: ["/control-plane"]
    working_dir: /config/
    command:
    - "--storage.type=file"
    - "--storage.refreshrate=1s"
    - "--secrets.type=file"
    - "--secrets.refreshrate=1s"
    - "--files.type=file"
    - "--files.refreshrate=1s"
    - "--xds.port=8081"
    volumes:
    - ./gloo-config:/config/

  function-discovery:
    image: soloio/function-discovery:0.2.0
    entrypoint: ["/function-discovery"]
    working_dir: /config/
    environment:
      GRPC_TRACE: all
      DEBUG: "1"
    command:
    - "--storage.type=file"
    - "--storage.refreshrate=1m"
    - "--secrets.type=file"
    - "--secrets.refreshrate=1m"
    - "--files.type=file"
    - "--files.refreshrate=1m"
    volumes:
    - ./gloo-config:/config/`
)
