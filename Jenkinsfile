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
                 image "golang:1.13"
             }
         }
         steps {
            sh "make test"
            sh "make testacc"
         }
      }
   }
}

