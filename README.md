## 前言
>
> 学校的项目需要通过监控考生，实现反作弊功能。
> 之前用python通过face_recognition库实现了人脸识别，今天想试试go是否可行。
> golang的go-face包需要安装dlib，我是用的macbook m1开发，
> 在引用go-face库时候出些各种问题，诸如 clang: error: the clang compiler does not support '-march=native'，又或者是 jpeg_mem_loader.cc:3:10: fatal error: 'jpeglib.h' file not found 等等。搞了半天都没弄通，于是想在ubuntu中试一试。
>


### docker拉取ubuntu镜像
    docker pull ubuntu
    docker run -itd --name myface-demo -v $PWD/myface/go:/var/go ubuntu
    
### 安装dlib、golang、
    docker exec -it myface-demo /bin/bash
    // 更新源
    apt-get update
    // 安装dlib gcc库
    apt-get install libdlib-dev libopenblas-dev libatlas-base-dev libjpeg-turbo8-dev build-essential wget vim
    // 查看版本号
    gcc --version
    gcc (Ubuntu 9.3.0-17ubuntu1~20.04) 9.3.0
    Copyright (C) 2019 Free Software Foundation, Inc.
    This is free software; see the source for copying conditions.  There is NO
    warranty; not even for MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.
### 安装golang([下载地址](https://golang.google.cn/dl/))
    wget https://golang.google.cn/dl/go1.17.7.linux-arm64.tar.gz
    // 解压到/usr/local目录    
    tar -C /usr/local -xzf go1.17.7.linux-arm64.tar.gz
    // 添加到环境变量中
    vim ~/.bashrc
    export PATH=$PATH:/usr/local/go/bin
    source ~/.bashrc
    // 检查版本
    go version
    go version go1.17.7 linux/arm64
### 拉取go-face包
    cd /var/go/demo
    go mod init myface
    touch main.go
    // 设置golang代理
    echo "export GO111MODULE=on" >> ~/.bashrc
    echo "export GOPROXY=https://goproxy.cn" >> ~/.bashrc
    source ~/.bashrc
    // 拉取go-face包
    go get github.com/Kagami/go-face
#### 项目结构
    --models
    --main.go
    --go.mod
    --jay-zhou.jpeg
    --jay.jpeg
#### 代码 main.go
    package main
    
    import (
    "fmt"
    "github.com/Kagami/go-face"
    "log"
    "path/filepath"
    )
    
    func main() {
    fmt.Println("Facial Recognition System v0.01")
    
    rec, err := face.NewRecognizer("models")
    if err != nil {
        fmt.Println("Cannot INItialize recognizer")
        return
    }
    defer rec.Close()
    
    fmt.Println("Recognizer Initialized")
    
    avengersImage := filepath.Join("jay-zhou.jpeg")
    
    faces, err := rec.RecognizeFile(avengersImage)
    if err != nil {
        log.Fatalf("Can't recognize: %v", err)
    }
    fmt.Println("Number of Faces in Image: ", len(faces))
    
    var samples []face.Descriptor
    var avengers []int32
    for i, f := range faces {
        samples = append(samples, f.Descriptor)
        // Each face is unique on that image so Goes to its own category.
        avengers = append(avengers, int32(i))
    }
    // Name the categories, i.e. people on the image.
    labels := []string{
        "周杰伦",
    }
    // Pass samples to the recognizer.
    rec.SetSamples(samples, avengers)
    
    // Now let's try to classify some not yet known image.
    testTonyStark := filepath.Join("jay.jpeg")
    tonyStark, err := rec.RecognizeSingleFile(testTonyStark)
    if err != nil {
        log.Fatalf("Can't recognize: %v", err)
    }
    if tonyStark == nil {
        log.Fatalf("Not a single face on the image")
    }
    avengerID := rec.ClassifyThreshold(tonyStark.Descriptor, 0.3)
    if avengerID < 0 {
        fmt.Println(avengerID)
        log.Fatalf("Can't classify")
    }
    
    fmt.Println(avengerID)
    
    fmt.Println(labels[avengerID])
    
    }

#### 执行main.go
    go run main.go
    Facial Recognition System v0.01
    Recognizer Initialized
    Number of Faces in Image:  1
    0
    周杰伦


