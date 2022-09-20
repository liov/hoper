package xyz.hoper.test.decimal;

import java.math.BigDecimal;
import java.math.RoundingMode;

public class Decimal {
    public static void main(String[] args) {
        var a = new BigDecimal("1.235");
        var b = new BigDecimal("1.23");
        System.out.println(a.divide(b, RoundingMode.HALF_UP));
    }
}
