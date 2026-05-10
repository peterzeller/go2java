public class Main {

    private static String goFmt(Object o) {
        if (o instanceof java.util.List<?> l) return l.toString().replace(", ", " ");
        return String.valueOf(o);
    }

    public static java.util.List<Integer> takeMid(java.util.List<Integer> in) {
        return in.subList(1, 3);
    }

    public static void main(String[] args) {
        var a = new java.util.ArrayList<>(java.util.List.of(1, 2, 3, 4));
        a.add(5);
        var b = a.subList(1, 4);
        var c = takeMid(a);
        System.out.println(goFmt(a));
        System.out.println(goFmt(b));
        System.out.println(goFmt(c));
    }
}
