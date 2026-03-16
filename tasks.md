Coisas a implementar
- [ ] API em go
    - [X] proxy reverso
        quando envia pra (proxy /teste/teste2 => servidor /teste/teste2)
        inicialmente bate em qualquer endpoint que receber
        depois podemos detectar se esse endpoint existe ?
    - [X] recebe request e da o forward nela e devolve a response
- [ ] da um jeito de receber de forma dinamica a lista de servidores ne
    - le de um arquivo ? CLI ? 
    - a gente le a URL; healthcheck.
- [ ] algoritmo pra reveza as requests (natan)
    -- escalonar pra mandar 1x pra cada servidor
    -- goroutine

- [ ] 