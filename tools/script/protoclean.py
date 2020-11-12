import os
import sys, getopt


def clean(inputpath, outputpath):
    if not os.path.exists(outputpath):
        os.makedirs(outputpath)
        clean(inputpath, outputpath)
    for root, dirs, files in os.walk(inputpath):
        for _dir in dirs:
            path = os.path.join(root, _dir)
            if path.endswith("utils/proto"):
                continue
            clean(path, os.path.join(outputpath, _dir))
        for file in files:
            replace(os.path.join(root, file))


def replace(inputfile):
    if not inputfile.endswith(".proto"):
        return


def main(argv):
    inputpath = ''
    outputpath = ''
    try:
        opts, args = getopt.getopt(argv, "hi:o:", ["ipath=", "opath="])
    except getopt.GetoptError:
        print('test.py -i <inputpath> -o <outputpath>')
        sys.exit(2)
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


if __name__ == "__main__":
    main(sys.argv[1:])
