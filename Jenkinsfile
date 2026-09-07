pipeline {
    agent none

   environment {
      SYSDIG_MONITOR_API_TOKEN = credentials('tech-marketing-token-monitor-lab')
      SYSDIG_SECURE_API_TOKEN = credentials('tech-marketing-token-secure-lab')
      GOCACHE="/tmp/go-build"
   }

   stages {
      stage('Check code') {
          agent {
              docker {
                  image "commitsar/commitsar"
              }
          }
          steps {
              warnError('Conventional Commits not being followed') {
                sh "commitsar"
              }
          }
      }
      stage('Tests') {
         agent {
             docker {
                 image "golang:1.27"
             }
         }
         steps {
            // Plain agent, no nix/just available here — inline what
            // 'just test'/'just testacc' actually run instead. Single-quoted
            // so ${TEST:-...} is expanded by sh, not Groovy, keeping the same
            // TEST/TESTARGS/TEST_SUITE override interface just/make had.
            sh './scripts/gofmtcheck.sh'
            sh 'go test ${TEST:-./...} -tags=unit -timeout=30s -parallel=4'
            sh 'CGO_ENABLED=1 TF_ACC=1 go test ${TEST:-./...} -v ${TESTARGS:-} -tags=${TEST_SUITE:-tf_acc_sysdig_monitor,tf_acc_sysdig_secure} -timeout 120m -race -parallel=1'
         }
      }
   }
}

