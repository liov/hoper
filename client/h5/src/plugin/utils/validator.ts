const Validator = {
  PhoneReg: new RegExp("^1[0-9]{10}$"),
  phone(input: string): boolean {
    return this.PhoneReg.test(input);
  },
  EmailReg:
    /^([a-zA-Z0-9]+[_.]?)*[a-zA-Z0-9]+@([a-zA-Z0-9]+[_.]?)*[a-zA-Z0-9]+.[a-zA-Z]{2,3}$/,
  mail(input: string): boolean {
    return this.EmailReg.test(input);
  },
};

export default Validator;
