# Stress Test CLI

### Como usar

Primeiro é necessário fazer o build da imagem docker rodando o seguinte comando na pasta raíz do projeto:
```
docker build -t stresstest .
```

Após finalizar o build da imagem, é possível iniciar os testes de carga com o seguinte comando:
```
docker run stresstest --url=http://google.com --requests=30 --concurrency=10
```

As variáveis URL, Requests e Concurrency são obrigatórias.
Ao final dos testes um relatório será informado contendo:
- Tempo total de execução
- Número total de requests
- Número de requests com sucesso
- Número de falhas com respectivos status HTTP