# Inicializar o init do projeto
```bash
go mod init github.com/fabiosoaresv/go_api
```

# Instalar lib
```bash
# Lib para requests
go get github.com/go-resty/resty/v2
# Lib para montar rotas
go get github.com/go-chi/chi/v5
# Framework
# go get -u github.com/gin-gonic/gin
go install github.com/traefik/yaegi/cmd/yaegi@latest
yaegi
```

# Debug
```bash
go install github.com/go-delve/delve/cmd/dlv@latest
dlv debug
break main.main
continue -> vai até o breakpoint
next -> pula linha
print variavel -> imprime a variável
```

# Rodar o projeto
```bash
go run cmd/server/main.go
```

# Config
```bash
# No arquivo go.mod você precisa apontar quem é o seu projeto, se for local, só definir o mesmo nome
module go_api
```


# Exemplos de códigos
## Criando função
```go
  func sum(n1 int, n2 int) {
    return n1 + n2
  }
```

## Criando struct
```go
  type User struct {
    Name string
    Age int
  }
```

## Declarar ponteiros
```go
user := &User{Name: "Fabio", Age: 28}

var user *User
user = &User{Name: "Fabio", Age: 28}
```

## Declarar array
```go
list := []string{"Fabio", "Soares", "Venturelli"}

# Append
append(list, "Dev")

# Remove
index := 1
list = append(list[:index], list[index+1:]...)
```

## Declarar hash
TODO: terminar de estudar aqui
```go
idades := map[string]int{
    "Fabio": 28,
    "João":  30,
}
```

## Goroutine
Uma goroutine é uma função ou método que é executado de forma concorrente no Go. Elas são gerenciadas pelo Go runtime, que distribui as goroutines entre as threads do sistema

## Go Scheduler
O Go scheduler é responsável por gerenciar as goroutines, escalonando e distribuindo sua execução entre as threads do sistema (ou threads do próprio runtime) de forma eficiente. Ele utiliza o número de CPUs disponíveis, determinado pela configuração de runtime.GOMAXPROCS(), para gerenciar o paralelismo de forma otimizada.

## Channel
Estrutura para comunicação segura entre goroutines, permitindo o envio e recebimento de dados.
- Você tem concorrência (várias goroutines rodando em paralelo) e precisa compartilhar dados entre elas com segurança.

## Race condition
Ocorre quando múltiplas goroutines acessam e modificam recursos compartilhados simultaneamente, resultando em comportamentos imprevisíveis.

## Mutex
Mutual Exclusion - é uma trava para garantir que apenas uma goroutine acesse um recurso compartilhado por vez, garantindo que uma goroutine termine antes de outra começar para não dar conflito nos resultados

## Deadlock
Todas as goroutines estão bloqueadas esperando por algo que nunca virá.

### Principais causas:
| Causa                               | Exemplo                         | Solução                                   |
| ----------------------------------- | ------------------------------- | ----------------------------------------- |
| Envio para canal sem receptor       | `ch <- 1` e ninguém lendo       | Adicionar goroutine consumidora ou buffer |
| Leitura de canal fechado            | `x := <-ch` com `ch` já fechado | Verificar se canal deve ser fechado       |
| `main` termina antes das goroutines | `go f()` e `main` termina       | Usar `WaitGroup`                          |
| `range` em canal nunca fechado      | `for x := range ch`             | Certificar que canal será fechado         |

Exemplo com solução
```go
package main

import (
	"fmt"
  // precisamos importar a lib sync
	"sync"
)

// passamos como parâmetro o ponteiro do WaitGroup na função
func worker(id int, jobs <-chan int, results chan<- int, wg *sync.WaitGroup) {
  // indica que a goroutine terminou. Por convenção, sempre que usamos WaitGroup, chamamos Done com defer no início da função.	defer wg.Done()
	for j := range jobs {
		results <- j * 2
	}
}

func main() {
	jobs := make(chan int, 5)
	results := make(chan int, 3)
  // declaramos uma variável de wait group
	var wg sync.WaitGroup

	for w := 1; w <= 3; w++ {
    // incrementa o contador para sinalizar uma nova goroutine pendente
		wg.Add(1)
		go worker(w, jobs, results, &wg)
	}

	for j := 1; j <= 5; j++ {
		jobs <- j
	}

   // sempre devemos fechar canais de escrita quando não serão mais usados
	close(jobs)

	go func() {
		// como o canal `results` será lido, precisamos esperar todas as goroutines terminarem com `wg.Wait()` e só então fechamos o canal com segurança
		wg.Wait()
		close(results)
	}()

	for a := 1; a <= 5; a++ {
		fmt.Println(<-results)
	}
}
```

## Wait Group
| Conceito          | `sync.Mutex`                               | `sync.WaitGroup`                              |
| ----------------- | ------------------------------------------ | --------------------------------------------- |
| **Serve para**    | Proteger acesso a **dados compartilhados** | **Esperar** goroutines terminarem             |
| **Evita**         | Condições de corrida (race conditions)     | O programa seguir antes de todas terminarem   |
| **Exemplo comum** | Incrementar uma variável com segurança     | Esperar todas as goroutines salvarem no banco |
| **Métodos**       | `.Lock()` e `.Unlock()`                    | `.Add()`, `.Done()`, `.Wait()`                |

## Exemplo
```go
var wg sync.WaitGroup

wg.Add(2) // vamos esperar 2 goroutines

go func() {
    // defer do ingles de "adiar" / "postergar"
    defer wg.Done() // sinaliza que terminou
    fmt.Println("Goroutine 1")
}()

go func() {
    defer wg.Done() // sinaliza que terminou
    fmt.Println("Goroutine 2")
}()

wg.Wait() // espera as duas goroutines finalizarem
fmt.Println("Todas as goroutines terminaram")

```

## Exemplo de erro de parser
```go
package main

import (
	"fmt"
	"net/http"

	"github.com/golang-jwt/jwt/v4"
)

func jwtMiddleware(w http.ResponseWriter, r *http.Request) {
	tokenString := r.Header.Get("Authorization")
  // aqui tava passando token, _ ignorando o erro, passei a receber o erro e interpretar embaixo
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte("minha_chave"), nil
	})

  // aqui o erro não tinha sido declarado e ele tentava acessar a variável token.Valid e quebrava
	if err != nil {
		http.Error(w, "Token inválido", http.StatusUnauthorized)
		return
	} else if token.Valid {
		fmt.Fprintln(w, "Token válido")
	} else {
		http.Error(w, "Token inválido", http.StatusUnauthorized)
	}
}

func main() {
	http.HandleFunc("/", jwtMiddleware)
	http.ListenAndServe(":8080", nil)
}
```

## Exemplo de recursividade infinita
```go
package main

import "fmt"

func countdown(n int) {
  // aqui tava sem o if, ele ia ficar eternamente
	if n <= 0 {
		fmt.Println("Countdown finished!")
		return
	}

	fmt.Println(n)
	countdown(n - 1)
}

func main() {
	countdown(10)
}
```

## Exemplo de tratativa de erro
```go
package main

import (
	"fmt"
	"os"
)

func main() {
	f, err := os.Create("output.txt")
	if err != nil {
		fmt.Println("Erro ao criar o arquivo:", err)
		return
	}
	defer f.Close()

	_, err = f.WriteString("Olá, mundo!\n")

	if err != nil {
		fmt.Println("Erro ao escrever no arquivo:", err)
		return
	}
	fmt.Println("Arquivo criado com sucesso.")
}
```

## Exemplo Race Condition
```go
package main

import (
	"fmt"
	"sync"
)

func main() {
	var counter int
	var wg sync.WaitGroup
	var mu sync.Mutex

	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func() {
      // Usamos o Mutex para evitar condição de corrida ao acessar a variável compartilhada
			mu.Lock() // bloqueia o acesso para garantir exclusividade
			counter++  // região crítica: leitura/escrita da variável compartilhada
			mu.Unlock() // libera o acesso para outras goroutines
			wg.Done()
		}()
	}

	wg.Wait()
	fmt.Println("Counter:", counter)
}
```

## Marmota (_)
Utilizado para ignorar valores que não vou usar, exemplo
```go
res, _ := dividir(10, 2) // ignora o erro
```

## Testar código
```go
package main

import (
    "fmt"
    "sync"
)

var counter int
var mu sync.Mutex

func increment(wg *sync.WaitGroup) {
    defer wg.Done()
    for i := 0; i < 1000; i++ {
        mu.Lock()
        counter++
        mu.Unlock()
    }
}

func main() {
    var wg sync.WaitGroup

    // Lançando 10 goroutines para incrementar o contador
    for i := 0; i < 10; i++ {
        wg.Add(1)
        go increment(&wg)
    }

    // Aguarda todas as goroutines terminarem
    wg.Wait()

    fmt.Println("Valor final do contador:", counter)
}

go run main.go
```



  func sum(n1 int, n2 int) (int, string){
    return n1 + n2, "fabio"
  }
