FROM devopsfaith/krakend:2.0.6 as builder
ARG ENV=dev

COPY krakend.json .
COPY config .

## Save temporary file to /tmp to avoid permission errors
RUN FC_ENABLE=1 \
    FC_OUT=/tmp/krakend.json \
    FC_PARTIALS="/etc/krakend/partials" \
    FC_SETTINGS="/etc/krakend/settings/$ENV" \
    FC_TEMPLATES="/etc/krakend/templates" \
    krakend check -d -t -c krakend.json

# The linting needs the final krakend.json file
RUN krakend check -c /tmp/krakend.json --lint

FROM devopsfaith/krakend:2.0.6
COPY --from=builder --chown=krakend /tmp/krakend.json .
# Uncomment with Enterprise image:
# COPY LICENSE /etc/krakend/LICENSE