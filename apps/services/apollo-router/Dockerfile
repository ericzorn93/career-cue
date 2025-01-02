FROM ghcr.io/apollographql/router:v1.59.0

ENV HOST=0.0.0.0
ENV PORT=4000
EXPOSE 4000

COPY ./router-prod.yaml /router.yaml
COPY ./supergraph-prod.graphql /supergraph.graphql

CMD ["--config", "/router.yaml", "--supergraph", "/supergraph.graphql"]
