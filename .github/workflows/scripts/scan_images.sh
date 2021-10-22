#!/bin/bash
# Generated by internal/ci/ci_tool.cue; do not edit
if [[ -z "${TAG}" ]]; then
  echo "TAG isn't set"
  exit 1
fi

if [[ -z "${REDHAT_TOKEN}" ]]; then
  echo "REDHAT_TOKEN isn't set"
  exit 1
fi

echo "::group::Scan quay.io/rh-marketplace/redhat-marketplace-operator"
id=$(curl -X GET "https://catalog.redhat.com/api/containers/v1/projects/certification/pid/ospid-c93f69b6-cb04-437b-89d6-e5220ce643cd" -H  "accept: application/json" -H  "X-API-KEY: $REDHAT_TOKEN" | jq -r '._id')
digest=$(skopeo --override-os=linux inspect docker://quay.io/rh-marketplace/redhat-marketplace-operator:$TAG | jq -r '.Digest')
curl -X POST "https://catalog.redhat.com/api/containers/v1/projects/certification/id/$id/requests/scans" \
--header 'Content-Type: application/json' \
--header "X-API-KEY: $REDHAT_TOKEN" \
--data-raw "{\"pull_spec\": \"quay.io/rh-marketplace/redhat-marketplace-operator@$digest\",\"tag\": \"$TAG\"}"
echo "::endgroup::"
echo "::group::Scan quay.io/rh-marketplace/redhat-marketplace-metric-state"
id=$(curl -X GET "https://catalog.redhat.com/api/containers/v1/projects/certification/pid/ospid-9b9b0dbe-7adc-448e-9385-a556714a09c4" -H  "accept: application/json" -H  "X-API-KEY: $REDHAT_TOKEN" | jq -r '._id')
digest=$(skopeo --override-os=linux inspect docker://quay.io/rh-marketplace/redhat-marketplace-metric-state:$TAG | jq -r '.Digest')
curl -X POST "https://catalog.redhat.com/api/containers/v1/projects/certification/id/$id/requests/scans" \
--header 'Content-Type: application/json' \
--header "X-API-KEY: $REDHAT_TOKEN" \
--data-raw "{\"pull_spec\": \"quay.io/rh-marketplace/redhat-marketplace-metric-state@$digest\",\"tag\": \"$TAG\"}"
echo "::endgroup::"
echo "::group::Scan quay.io/rh-marketplace/redhat-marketplace-reporter"
id=$(curl -X GET "https://catalog.redhat.com/api/containers/v1/projects/certification/pid/ospid-faa0f295-e195-4bcc-a3fc-a4b97ada317e" -H  "accept: application/json" -H  "X-API-KEY: $REDHAT_TOKEN" | jq -r '._id')
digest=$(skopeo --override-os=linux inspect docker://quay.io/rh-marketplace/redhat-marketplace-reporter:$TAG | jq -r '.Digest')
curl -X POST "https://catalog.redhat.com/api/containers/v1/projects/certification/id/$id/requests/scans" \
--header 'Content-Type: application/json' \
--header "X-API-KEY: $REDHAT_TOKEN" \
--data-raw "{\"pull_spec\": \"quay.io/rh-marketplace/redhat-marketplace-reporter@$digest\",\"tag\": \"$TAG\"}"
echo "::endgroup::"
echo "::group::Scan quay.io/rh-marketplace/redhat-marketplace-authcheck"
id=$(curl -X GET "https://catalog.redhat.com/api/containers/v1/projects/certification/pid/ospid-ffed416e-c18d-4b88-8660-f586a4792785" -H  "accept: application/json" -H  "X-API-KEY: $REDHAT_TOKEN" | jq -r '._id')
digest=$(skopeo --override-os=linux inspect docker://quay.io/rh-marketplace/redhat-marketplace-authcheck:$TAG | jq -r '.Digest')
curl -X POST "https://catalog.redhat.com/api/containers/v1/projects/certification/id/$id/requests/scans" \
--header 'Content-Type: application/json' \
--header "X-API-KEY: $REDHAT_TOKEN" \
--data-raw "{\"pull_spec\": \"quay.io/rh-marketplace/redhat-marketplace-authcheck@$digest\",\"tag\": \"$TAG\"}"
echo "::endgroup::"
echo "::group::Scan quay.io/rh-marketplace/redhat-marketplace-data-service"
id=$(curl -X GET "https://catalog.redhat.com/api/containers/v1/projects/certification/pid/ospid-61649f78d3e2f8d3bcfe30d5" -H  "accept: application/json" -H  "X-API-KEY: $REDHAT_TOKEN" | jq -r '._id')
digest=$(skopeo --override-os=linux inspect docker://quay.io/rh-marketplace/redhat-marketplace-data-service:$TAG | jq -r '.Digest')
curl -X POST "https://catalog.redhat.com/api/containers/v1/projects/certification/id/$id/requests/scans" \
--header 'Content-Type: application/json' \
--header "X-API-KEY: $REDHAT_TOKEN" \
--data-raw "{\"pull_spec\": \"quay.io/rh-marketplace/redhat-marketplace-data-service@$digest\",\"tag\": \"$TAG\"}"
echo "::endgroup::"