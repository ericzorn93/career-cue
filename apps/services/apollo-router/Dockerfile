FROM ghcr.io/apollographql/router:v1.59.0

ENV HOST=0.0.0.0
ENV PORT=4000
EXPOSE 4000

WORKDIR /dist

COPY ./apps/services/apollo-router/router-prod.yaml ./config/router.yaml
COPY ./apps/services/apollo-router/supergraph-prod.graphql ./config/supergraph.graphql

CMD ["--config", "/dist//config/router.yaml", "--supergraph", "/dist/config/supergraph.graphql"]
