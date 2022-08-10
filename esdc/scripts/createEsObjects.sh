#!/bin/sh
ELASTICSEARCH_HOST="http://es-container:9200"
#ELASTICSEARCH_HOST="http://localhost:9200"

escurlwithretry() {
    max=150
    for i in $(seq 2 $max); do
        ret=$(escurl "$@")
        if [ "$ret" -gt 0 ]; then
            echo "[ERROR] Curl failed with exit code : $ret"
            sleep 1s
        else
            echo "[INFO] curl success"
            break
        fi
    done
}

escurl() {
    ret=0
    curl -H "Content-Type: application/json" -k "$@" >/dev/null 2>&1 || ret=$?
    if [ "$ret" -gt 0 ]; then
        echo "${ret}"
    else
        echo "0"
    fi
}

for file in ilm-policy-*.json; do
    name="$(basename ${file} .json | sed 's,ilm-policy-,,g')-policy"
    echo "[INFO] Deploying ILM policy ${name}"
    escurlwithretry -X PUT -d@${file} "${ELASTICSEARCH_HOST}/_ilm/policy/${name}"
done

for file in index-template-*.json; do
    name="$(basename ${file} .json | sed 's,index-template-,,g')"
    echo "[INFO] Deploying index template ${name}-template"
    escurlwithretry -X PUT -d@${file} "${ELASTICSEARCH_HOST}/_index_template/${name}-template"
    escurlwithretry -X PUT "${ELASTICSEARCH_HOST}/_data_stream/${name}"
done

for file in rollup-job-*.json; do
    name="$(basename ${file} .json | sed 's,rollup-job-,,g')-job"
    echo "[INFO] Deploying rollup job ${name}"
    escurlwithretry -X PUT -d@${file} "${ELASTICSEARCH_HOST}/_rollup/job/${name}"
    escurlwithretry -X POST "${ELASTICSEARCH_HOST}/_rollup/job/${name}/_start"
done

sleep 5s