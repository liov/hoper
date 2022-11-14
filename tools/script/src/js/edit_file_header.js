import  fs from 'fs'

for(let i = 1; i < 24;i++){
  let bin = fs.readFileSync(`D:\\F\\concat\\${i}.ts`);
  console.log(bin.subarray(0,4))
  bin[0] = 0xff
  bin[1] = 0xff
  bin[2] = 0xff
  bin[3] = 0xff
  fs.writeFileSync(`D:\\F\\concat\\${i}.ts`,bin)

}
