job "gloo" {

  datacenters = [
    "dc1"]
  type = "service"

  update {
    max_parallel = 1
    min_healthy_time = "10s"
    healthy_deadline = "3m"
    auto_revert = false
    canary = 0
  }

  migrate {
    max_parallel = 1
    health_check = "checks"
    min_healthy_time = "10s"
    healthy_deadline = "5m"
  }

  group "gloo" {
    count = 1
    restart {
      attempts = 2
      interval = "30m"
      delay = "15s"
      mode = "fail"
    }
    ephemeral_disk {
      size = 300
    }

    # control plane
    task "control-plane" {
      env {
        DEBUG = "1"
      }
      driver = "docker"
      config {
        image = "soloio/control-plane:0.2.0"
        port_map {
          xds = 8081
        }
        args = [
          "--storage.type=consul",
          "--storage.refreshrate=1m",
          "--secrets.type=vault",
          "--secrets.refreshrate=1m",
          "--files.type=consul",
          "--files.refreshrate=1m",
          "--xds.port=${NOMAD_PORT_xds}",
          "--consul.address=${attr.driver.docker.bridge_ip}:8500",
          "--consul.scheme=http",
          "--vault.addr=http://${attr.driver.docker.bridge_ip}:8200",
          "--vault.token=${VAULT_TOKEN}",
        ]
      }
      resources {
        cpu = 500
        memory = 256
        network {
          mbits = 10
          port "xds" {}
        }
      }
      service {
        name = "control-plane"
        tags = [
          "gloo"]
        port = "xds"
        check {
          name = "alive"
          type = "tcp"
          interval = "10s"
          timeout = "2s"
        }
      }
      vault {
        change_mode = "restart"
        policies = [
          "gloo"]
      }
    }

    # ingress
    task "ingress" {

      driver = "docker"
      config {
        image = "soloio/envoy:v0.1.6-127"
        port_map {
          http = 8080
          https = 8443
          admin = 19000
        }
        command = "/bin/bash"
        args = [
          "-c",
          <<EOF
cat > ${NOMAD_TASK_DIR}/envoy.yaml <<CONFIG_END; envoy -c /local/envoy.yaml --v2-config-only
node:
  cluster: ingress
  id: ${NOMAD_ALLOC_ID}

static_resources:
  clusters:

  - name: xds_cluster
    connect_timeout: 5.000s
    hosts:
    - socket_address:
        address: ${NOMAD_IP_control_plane_xds}
        port_value: ${NOMAD_PORT_control_plane_xds}
    http2_protocol_options: {}
    type: STATIC

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
  access_log_path: /dev/null
  address:
    socket_address:
      address: 0.0.0.0
      port_value: 19000
CONFIG_END

EOF
        ]
      }
      resources {
        cpu = 500
        memory = 256
        network {
          mbits = 10
          port "http" {}
          port "https" {}
          port "admin" {}
        }
      }
      service {
        name = "ingress"
        tags = [
          "gloo"]
        port = "admin"
        check {
          name = "alive"
          type = "tcp"
          interval = "10s"
          timeout = "5s"
        }
      }
    }

    # function-discovery
    task "function-discovery" {

      env {
        DEBUG = "1"
      }

      driver = "docker"
      config {
        image = "soloio/function-discovery:0.2.0"
        args = [
          "--storage.type=consul",
          "--storage.refreshrate=1m",
          "--secrets.type=vault",
          "--secrets.refreshrate=1m",
          "--files.type=consul",
          "--files.refreshrate=1m",
          "--consul.address=${attr.driver.docker.bridge_ip}:8500",
          "--consul.scheme=http",
          "--vault.addr=http://${attr.driver.docker.bridge_ip}:8200",
          "--vault.token=${VAULT_TOKEN}",
        ]
      }
      resources {
        cpu = 500
        memory = 256
      }
      vault {
        change_mode = "restart"
        policies = [
          "gloo"]
      }
    }



  }

}
