public class Main {

    private static String goFmt(Object o) {
        if (o instanceof java.util.List<?> l) return l.toString().replace(", ", " ");
        return String.valueOf(o);
    }

    public static void main(String[] args) {
        System.out.println(goFmt("hello from go"));
    }
}
