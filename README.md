# Mutant Challenge API
Esta aplicacion ayuda a determinar si un humano es mutante basándose en su secuencia de ADN.

La API contara con dos recursos:

 - Mutant
 - Stat

## Consumo API
El API expuesto podra ser consumido para realizar validaciones de ADN, solo las cadenas de ADN que sean validas (sin importar si son humanas o mutantes) seran almacenadas en la base de datos DynamoDB. La API recibirá como parámetro un array de Strings que representan cada fila de una tabla de (NxN) con la secuencia del ADN. Las letras de los Strings solo pueden ser: (A,T,C,G), las cuales representa cada base nitrogenada del ADN.

![Table Image](https://github.com/RiAlGueFor/mutant-challenge/blob/main/img/ejemplos_DNA.png)

El API sera expuesta desde AWS usando una funcion lambda de la siguiente manera:

![Table Image](https://github.com/RiAlGueFor/mutant-challenge/blob/main/img/api-gateway.png)

### Recurso Mutants

Este recurso sera expuesto con el metodo POST y validara si la secuencia de ADN pertenece a un humano o a un mutante. En caso de verificar un mutante, el metodo  devolverá un HTTP 200-OK, en caso contrario un 403-Forbidden

cURL de consumo: 

curl --location --request POST 'https://vcg41iv3i4.execute-api.sa-east-1.amazonaws.com/staging/mutant' \
--header 'Content-Type: application/json' \
--data-raw '{
    "dna": [
        "ATGCGA",
        "CAGTAC",
        "TCAGGT",
        "AGCGGG",
        "ATACAG",
        "TCACTG"
    ]
}'

Resultado Humano

![Table Image](https://github.com/RiAlGueFor/mutant-challenge/blob/main/img/mutant-POST-Humano.png)

Resultado Mutante

![Table Image](https://github.com/RiAlGueFor/mutant-challenge/blob/main/img/mutant-POST-Mutante.png)

Cadenas de ADN invalidas

![Table Image](https://github.com/RiAlGueFor/mutant-challenge/blob/main/img/mutant-POST-IncorrectSize.png)
![Table Image](https://github.com/RiAlGueFor/mutant-challenge/blob/main/img/mutant-POST-IncorrectChain.png)

Cuando las cadenas de ADN son válidas, sin importar que tras la validación sean mutantes o humanas, las cadenas de ADN serán almacenadas en la base de datos DynamoDB 

![Table Image](https://github.com/RiAlGueFor/mutant-challenge/blob/main/img/DynamoDB.png)

Este Recurso se encuentra definido en AWS API Gateway como 

![Table Image](https://github.com/RiAlGueFor/mutant-challenge/blob/main/img/api-gateway-POST%20-%20copia.png)

### Recurso Stats

Este recurso sera expuesto con el metodo GET y se encarga de retornar las estadísticas de las verificaciones de ADN

cURL de consumo: 
curl --location --request GET 'https://vcg41iv3i4.execute-api.sa-east-1.amazonaws.com/staging/stats'

![Table Image](https://github.com/RiAlGueFor/mutant-challenge/blob/main/img/stats-GET.png)

Este Recurso se encuentra definido en AWS API Gateway como 

![Table Image](https://github.com/RiAlGueFor/mutant-challenge/blob/main/img/api-gateway-GET.png)

## Cobertura


## Requisitos previos

Para seguir el ejemplo de este artículo, necesitará lo siguiente:

  - Un espacio de trabajo de Go configurado conforme al siguiente [enlace](https://go.dev/doc/install)
  - Una cuenta de AWS en la cual seran necesarios crear una [funcion Lambda](https://docs.aws.amazon.com/lambda/latest/dg/lambda-golang.html) y un [API Gateway](https://docs.aws.amazon.com/apigateway/latest/developerguide/apigateway-rest-api.html) para exponer un servicio REST.

## Compilar el paquete

El proyecto ha sido desarrollado en Go, inicialmente se requiere realizar los siguientes pasos para generar el ejecutable:

  - En una consola de comandos nos dirigimos a la ruta del codigo (src/go-lambda-mutants)
  - Una vez estamos ubicados en la ruta ejecutamos los siguiente comandos

```
go mod init [nombre_modulo]
go mod tidy
```

### En Linux and macOS

```
set GO111MODULE=on
go.exe get -u github.com/aws/aws-lambda-go/cmd/build-lambda-zip
```

### En Windows

Obtenemos la herramienta
```
set GO111MODULE=on
go.exe get -u github.com/aws/aws-lambda-go/cmd/build-lambda-zip
```
Use la herramiento desde su GOPATH (en algunos casos puede ser GOBIN). Si tienes la instalacion por defecto de Go, la herramienta esta en la ruta  %USERPROFILE%\Go\bin.

En cmd.exe:
```
set GOOS=linux
set GOARCH=amd64
set CGO_ENABLED=0
go build -o main main.go
%USERPROFILE%\Go\bin\build-lambda-zip.exe -o main.zip main
```

## Despliegue

Con el zip generado procedemos a cargarlo en una funcion Lamdba previamente creada.

![Table Image](https://github.com/RiAlGueFor/mutant-challenge/blob/main/img/lamba-function.png)

Y en el API Gateway nos aseguramos de (re-)desplegar el API cada vez que se realice algun cambio.

## Tecnologias

Este proyecto fue implementado usando Golang y usa las siguientes tecnologias:

    -go version go1.18.4 windows/amd64
    -AWS
      - DynamoDB
      - Lambda Functions
      - API Gateway
    -Git
    
