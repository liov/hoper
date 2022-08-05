import os
import re
import sys, getopt


def clean(inputpath, outputpath):
    print(inputpath, "[", outputpath, "]")
    if not os.path.exists(outputpath):
        os.makedirs(outputpath)
    files = os.listdir(inputpath)
    for file in files:
        path = os.path.join(inputpath, file)
        if os.path.isfile(path):
            replace(path, os.path.join(outputpath, file))
        else:
            clean(path, os.path.join(outputpath, file))


def replace(inputfile, outputfile):
    if not inputfile.endswith(".proto"):
        return
    f = open(inputfile, "r", encoding="utf-8")
    data = f.read()
    f.close()
    data = re.sub(r"import \"github.*\n", "", data)
    data = re.sub(r"import \"protoc-gen-openapiv2.*\n", "", data)
    data = re.sub(r"import \"utils/proto/gogo/enum.proto.*\n", "", data)
    data = re.sub(r"import \"utils/proto/gogo/graphql.proto.*\n", "", data)
    data = re.sub(r"import \"google/api/annotations.proto.*\n", "", data)
    data = re.sub(r"\[\([\w\s.)=\"/:\\;\-(\'{},一-龥，]*]", "", data)
    data = re.sub(r"option \([\w\s.)/\[\]=\":*@\-('{},一-龥]*;", "", data)
    f = open(outputfile, "w+", encoding="utf-8")
    f.write(data)
    f.close()


def main(argv):
    inputpath = '../../proto'
    outputpath = '../../proto_std'

    opts, args = getopt.getopt(argv, "hi:o:", ["ipath=", "opath="])

    for opt, arg in opts:
        if opt == '-h':
            print('test.py -i <inputpath> -o <outputpath>')
            sys.exit()
        elif opt in ("-i", "--ipath"):
            inputpath = arg
        elif opt in ("-o", "--opath"):
            outputpath = arg
    print('输入的路径为：', inputpath)
    print('输出的路径为：', outputpath)
    clean(inputpath, outputpath)


if __name__ == "__main__":
    main(sys.argv[1:])
