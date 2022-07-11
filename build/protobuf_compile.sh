#!/bin/bash

os=$(uname -s)
GOOS=""
case ${os} in
    "Linux" ) 
        GOOS="linux"
    ;;
    "Darwin" ) 
        GOOS="darwin"
    ;;
    * ) 
        echo -e "[\033[34mFatal\033[0m]: 暂不支持的系统类型[${os}] "
        exit 255
    ;;
esac

arch=$(uname -m)
GOARCH=""
case ${arch} in
    "x86_64" ) 
        GOARCH="amd64"
    ;;
    * ) 
        echo -e "[\033[34mFatal\033[0m]: 暂不支持的 cpu 架构[${arch}] "
        exit 255
    ;;
esac

root_dir=$(pwd)

for name in `ls ${root_dir}/app`
do
    app_dir=${root_dir}/app/${name}
    api_dir=${app_dir}/api
    if [ -d ${api_dir} ]; then
        ok=false
        for name in `ls ${api_dir}`
        do
            if [ "${file##*.}"x = "proto"x ]; then
                ok=true
                break
            fi
        done

        if [ ${ok}=true ]; then
            pb_dir=${app_dir}/pb
            if [ -d ${pb_dir} ]; then
                rm -rf ${pb_dir}
            fi 
            mkdir -p ${pb_dir}
            cd ${api_dir}

            for name in `ls ${api_dir}`
            do
                file_name=${api_dir}/${name}
                ${root_dir}/tools/protoc/${GOOS}_${GOARCH}/bin/protoc --proto_path="" -I . --go_out=${pb_dir} --go-grpc_out=${pb_dir} ${name}
                code=$(echo $?)
                if [ $code = 0 ]; then
                    echo -e "编译文件: ${file_name} => [\033[31m成功\033[0m] "
                else
                    echo -e "[\033[34mFatal\033[0m]: 编译文件: [${file_name}] => [\033[34m失败\033[0m] "
                    echo -e "\t <<<<<<<<<<<< 编译过程意外退出，已终止  <<<<<<<<<<<<"
                    exit
                fi
            done
        fi
    fi
done
