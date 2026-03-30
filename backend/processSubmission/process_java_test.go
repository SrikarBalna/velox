package processSubmission_test

import (
	"testing"
)

func TestProcessSubmission_Java(t *testing.T) {
	tests := []ExecutionTestCase{
		{
			Name: "Java_Accepted",
			SourceCode: `import java.util.Scanner;
public class Main {
    public static void main(String[] args) {
        Scanner sc = new Scanner(System.in);
        if (sc.hasNextInt()) {
            int a = sc.nextInt();
            int b = sc.nextInt();
            System.out.print(a + b);
        }
    }
}`,
			ExpectedState: "Accepted",
		},
		{
			Name: "Java_Wrong_Answer",
			SourceCode: `import java.util.Scanner;
public class Main {
    public static void main(String[] args) {
        System.out.print("1000");
    }
}`,
			ExpectedState: "Wrong Answer",
		},
		{
			Name: "Java_Compile_Error_NoMainClass",
			SourceCode: `public class Solution {
    public static void main(String[] args) { System.out.print("Hi"); }
}`, // Expected Main class because Java compilation writes to Main.java
			ExpectedState: "Compile Error",
		},
		{
			Name: "Java_Compile_Error_Syntax",
			SourceCode: `public class Main {
    public static void main(String[] args) { int a = "string"; }
}`,
			ExpectedState: "Compile Error",
		},
		{
			Name: "Java_Runtime_Error_Exception",
			SourceCode: `public class Main {
    public static void main(String[] args) { 
		int[] arr = new int[5];
		System.out.print(arr[10]); // ArrayIndexOutOfBoundsException
	}
}`,
			ExpectedState: "Runtime Error",
		},
		{
			Name: "Java_Time_Limit_Exceeded",
			SourceCode: `public class Main {
    public static void main(String[] args) throws Exception { 
		Thread.sleep(5000);
	}
}`,
			TimeLimitMs: 1500, // Slightly higher for Java warmup
			ExpectedState: "Time Limit Exceeded",
		},
		// Note: Java MLE is notoriously tricky under Rusage due to JVM overhead. 
		// We allocate safely and verify it restricts JVM heap growth if limits are enforced properly by OS.
		{
			Name: "Java_Memory_Limit_Exceeded",
			SourceCode: `import java.util.ArrayList;
public class Main {
    public static void main(String[] args) {
        ArrayList<byte[]> list = new ArrayList<>();
        while(true) {
            list.add(new byte[1024 * 1024]);
        }
    }
}`,
			MemoryLimitKb: 150 * 1024, // JVM baseline is high, setting limit to 150MB to trap it
			ExpectedState: "Memory Limit Exceeded", // Might evaluate to Runtime Error (OOM Exception caught by bash execution). Handled generically here via fail fast.
		},
	}

	tests = append(tests, GenerateVariations("Java_Bulk_Accepted", 
		`import java.util.Scanner;
		public class Main {
			public static void main(String[] args) {
				Scanner sc = new Scanner(System.in);
				if (sc.hasNextInt()) {
					System.out.print(sc.nextInt() + sc.nextInt());
				}
			}
		}`, 
		"Accepted", 15)...)

	runLanguageTests(t, "java", tests)
}
