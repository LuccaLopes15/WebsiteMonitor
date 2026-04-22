package main

import (
	"bufio"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

func main() {
	fmt.Println("Bem vindo ao GoMonitor!")

	for {
		exibirMenu()
		comando := lerComando()

		switch comando {
		case 1:
			iniciarMonitoramento()
		case 2:
			verificarUrlsMonitoradas()
		case 3:
			novaUrl := lerNovaUrl()
			adicionarUrlNoMonitoramento(novaUrl)
		case 5:
			encerrar()
		}
	}
}

func exibirMenu() {
	fmt.Println("Escolha uma das opções:")
	fmt.Println("1 - Iniciar monitoramento")
	fmt.Println("2 - Verificar urls monitoradas")
	fmt.Println("3 - Adicionar url no monitoramento")
	fmt.Println("4 - Exibir logs")
	fmt.Println("5 - Encerrar")
}

func lerComando() int {
	var comando int
	fmt.Scan(&comando)
	return comando
}

func lerNovaUrl() string {
	fmt.Println("Digite a nova url que deseja monitorar: ")
	var novaUrl string
	fmt.Scan(&novaUrl)
	return novaUrl
}

func iniciarMonitoramento() {
	fmt.Println("Iniciando monitoramento...")

	sites, deuErro := lerSitesParaMonitorar()

	if deuErro {
		return
	}

	tamanho := len(sites)

	if tamanho <= 0 {
		fmt.Println("Não há sites para monitorar!")
		return
	}

	for _, site := range sites {
		res, err := http.Get(site)

		deuErro = gerouErro(err, "Erro ao tentar acessar site "+site)

		if deuErro {
			escreveNoLog(site, false)
			continue
		}

		if res.StatusCode == 200 {
			fmt.Printf("Site: %s foi carregado com sucesso!\n", site)
			escreveNoLog(site, true)
		} else {
			fmt.Printf("Site: %s está com erro. Status: %d\n", site, res.StatusCode)
			escreveNoLog(site, false)
		}
	}
}

func escreveNoLog(site string, status bool) {
	// os.OpenFile permite configurar o arquivo
	// O_RDWR: Leitura e Escrita
	// O_CREATE: Cria o arquivo se não existir
	// O_APPEND: Adiciona ao final do arquivo em vez de sobrescrever
	arquivo, err := os.OpenFile("log.txt", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)

	deuErro := gerouErro(err, "Erro ao abrir arquivo de log")

	if deuErro {
		return
	}

	defer arquivo.Close()

	// Usamos o pacote time para pegar a data atual formatada
	horario := time.Now().Format("02/01/2006 15:04:05")

	// Escreve no arquivo: Data - Site - online: true/false
	arquivo.WriteString(horario + " - " + site + " - online: " + strconv.FormatBool(status) + "\n")
}

func lerSitesParaMonitorar() ([]string, bool) {
	var sites []string

	fmt.Println("Carregando sites...")

	arquivo, err := os.Open("sites.txt")

	result := gerouErro(err, "Erro ao abrir arquivo")

	if result {
		return sites, true
	}

	defer arquivo.Close()

	scanner := bufio.NewScanner(arquivo)

	tot := 0
	for scanner.Scan() {
		linha := scanner.Text()
		linha = strings.TrimSpace(linha)

		if linha != "" {
			sites = append(sites, linha)
			tot = tot + 1
			fmt.Println("")
			fmt.Printf("Site %d carregado - %s", tot, linha)
		}
	}

	result = gerouErro(scanner.Err(), "Erro ao ler sites do arquivo")

	if result {
		return sites, true
	}

	fmt.Println("Sites carregados")

	return sites, false
}

func verificarUrlsMonitoradas() {
	sites, deuErro := lerSitesParaMonitorar()

	if deuErro {
		return
	}

	fmt.Println("Urls monitoradas: ")
	for _, site := range sites {
		fmt.Println(site)
	}
}

func adicionarUrlNoMonitoramento(novaUrl string) {
	arquivo, err := os.OpenFile("sites.txt", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)

	deuErro := gerouErro(err, "Erro ao abrir arquivo de sites")

	if deuErro {
		return
	}

	defer arquivo.Close()

	arquivo.WriteString("\n" + novaUrl)
}

func encerrar() {
	os.Exit(0)
}

func gerouErro(err error, message string) bool {
	if err != nil {
		fmt.Println(message, err)
		return true
	}
	return false
}
