"use client";

import React, { useState } from "react";

const API_URL = process.env.NEXT_PUBLIC_API_URL || "http://localhost:8080";

const LANGUAGE_TEMPLATES: Record<string, string> = {
  python: `import sys\n\ndef main():\n    # Read all lines from standard input\n    input_data = sys.stdin.read().splitlines()\n    if not input_data:\n        return\n    \n    # Example: Print back the first line\n    print(input_data[0])\n\nif __name__ == '__main__':\n    main()`,
  javascript: `const fs = require('fs');\n\nfunction main() {\n    // Read all input from standard input\n    const input = fs.readFileSync('/dev/stdin', 'utf-8').trim().split('\\n');\n    if (input.length === 0 || input[0] === '') return;\n\n    // Example: Print back the first line\n    console.log(input[0]);\n}\n\nmain();`,
  go: `package main\n\nimport (\n\t"bufio"\n\t"fmt"\n\t"os"\n)\n\nfunc main() {\n\tscanner := bufio.NewScanner(os.Stdin)\n\t// Example: Read first line and print it back\n\tif scanner.Scan() {\n\t\tfmt.Println(scanner.Text())\n\t}\n}`,
  cpp: `#include <iostream>\n#include <string>\n\nusing namespace std;\n\nint main() {\n    string line;\n    // Example: Read first line and print it back\n    if (getline(cin, line)) {\n        cout << line << endl;\n    }\n    return 0;\n}`,
  c: `#include <stdio.h>\n\nint main() {\n    char line[1024];\n    // Example: Read first line and print it back\n    if (fgets(line, sizeof(line), stdin)) {\n        printf("%s", line);\n    }\n    return 0;\n}`,
  java: `import java.util.Scanner;\n\npublic class Main {\n    public static void main(String[] args) {\n        Scanner scanner = new Scanner(System.in);\n        // Example: Read first line and print it back\n        if (scanner.hasNextLine()) {\n            System.out.println(scanner.nextLine());\n        }\n        scanner.close();\n    }\n}`,
};

export default function PlaygroundPage() {
  const [apiKey, setApiKey] = useState("");
  const [language, setLanguage] = useState("python");
  const [code, setCode] = useState(LANGUAGE_TEMPLATES.python);
  const [testCases, setTestCases] = useState([
    { id: 1, input: "", expectedOutput: "" },
  ]);
  const [timeLimit, setTimeLimit] = useState(2000);
  const [memoryLimit, setMemoryLimit] = useState(256000);

  const [submissionId, setSubmissionId] = useState<string | null>(null);
  const [statusLog, setStatusLog] = useState<any[]>([]);
  const [isSubmitting, setIsSubmitting] = useState(false);
  const [isChecking, setIsChecking] = useState(false);

  const logsContainerRef = React.useRef<HTMLDivElement>(null);

  React.useEffect(() => {
    if (logsContainerRef.current) {
      logsContainerRef.current.scrollTop =
        logsContainerRef.current.scrollHeight;
    }
  }, [statusLog]);

  const appendLog = (
    type: "info" | "success" | "error",
    message: string,
    data?: any,
  ) => {
    setStatusLog((prev) => [
      ...prev,
      { type, message, data, timestamp: new Date() },
    ]);
  };

  const handleSubmit = async () => {
    if (!apiKey) {
      appendLog("error", "API Key is required");
      return;
    }

    setIsSubmitting(true);
    appendLog("info", "Submitting code...");

    try {
      const res = await fetch(`${API_URL}/submit`, {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
          Authorization: `Bearer ${apiKey}`,
        },
        body: JSON.stringify({
          language,
          source_code: code,
          test_cases: testCases.map((tc) => ({
            test_case_id: tc.id,
            input: tc.input,
            expected_output: tc.expectedOutput,
          })),
          time_limit_ms: timeLimit,
          memory_limit_kb: memoryLimit,
        }),
      });

      const data = await res.json().catch(() => null);

      if (!res.ok) {
        appendLog(
          "error",
          `Submit failed: ${res.status || res.statusText}`,
          data,
        );
        setIsSubmitting(false);
        return;
      }

      if (data?.submission_id) {
        setSubmissionId(data.submission_id);
      }
      appendLog("success", "Submission accepted!", data);
    } catch (err: any) {
      appendLog("error", `Error submitting code: ${err.message}`);
    } finally {
      setIsSubmitting(false);
    }
  };

  const handleCheckStatus = async () => {
    if (!apiKey) {
      appendLog("error", "API Key is required to check status");
      return;
    }
    if (!submissionId) {
      appendLog("error", "No submission ID available");
      return;
    }

    setIsChecking(true);
    appendLog("info", `Checking status for ${submissionId}...`);

    try {
      const res = await fetch(
        `${API_URL}/status?submission_id=${submissionId}`,
        {
          method: "GET",
          headers: {
            Authorization: `Bearer ${apiKey}`,
          },
        },
      );

      const data = await res.json().catch(() => null);

      if (!res.ok) {
        appendLog(
          "error",
          `Check status failed: ${res.status || res.statusText}`,
          data,
        );
        setIsChecking(false);
        return;
      }

      appendLog(
        "success",
        `Status retrieved (${data?.status || data?.overall_state || "unknown"})`,
        data,
      );
    } catch (err: any) {
      appendLog("error", `Error checking status: ${err.message}`);
    } finally {
      setIsChecking(false);
    }
  };

  return (
    <div className="flex flex-col gap-6 max-w-7xl mx-auto w-full pb-10">
      <div className="flex flex-col mb-4">
        <div className="inline-flex items-center gap-2 px-3 py-1 bg-surface border border-white/10 text-white text-xs font-bold uppercase tracking-wider w-fit mb-3">
          <span className="w-2 h-2 bg-primary"></span>
          API Playground
        </div>
        <h1 className="text-3xl font-black text-white tracking-tight">
          Execution Sandbox
        </h1>
        <p className="text-foreground/50 text-sm mt-2 max-w-2xl">
          Submit code execution requests and track outputs in real-time.
        </p>
      </div>

      <div className="grid grid-cols-1 lg:grid-cols-12 gap-6">
        {/* Request Configuration Form */}
        <div className="lg:col-span-6 flex flex-col gap-6">
          <div className="p-6 bg-surface border border-white/10">
            <div className="flex items-center justify-between mb-6 border-b border-white/10 pb-4">
              <h2 className="text-sm font-black text-white uppercase tracking-widest">
                Request Configuration
              </h2>
              <span className="px-2 py-1 text-[10px] bg-primary text-black font-black uppercase tracking-widest">
                POST /submit
              </span>
            </div>

            <div className="space-y-4">
              <div className="space-y-1">
                <label className="text-[10px] font-bold text-white/50 uppercase tracking-widest">
                  API Key (vlx_)
                </label>
                <input
                  type="text"
                  value={apiKey}
                  onChange={(e) => setApiKey(e.target.value)}
                  placeholder="Enter API key"
                  className="w-full bg-black border border-white/10 px-4 py-2 text-white outline-none focus:border-primary transition-colors font-mono text-sm"
                />
              </div>

              <div className="grid grid-cols-2 gap-4">
                <div className="space-y-1">
                  <label className="text-[10px] font-bold text-white/50 uppercase tracking-widest">
                    Runtime
                  </label>
                  <select
                    value={language}
                    onChange={(e) => {
                      const newLang = e.target.value;
                      setLanguage(newLang);
                      setCode(LANGUAGE_TEMPLATES[newLang] || "");
                    }}
                    className="w-full bg-black border border-white/10 px-4 py-2 text-white outline-none focus:border-primary transition-colors font-bold text-sm"
                  >
                    <option value="python">Python 3.10</option>
                    <option value="javascript">Node.js (V8)</option>
                    <option value="go">Go 1.21</option>
                    <option value="cpp">C++ (GCC)</option>
                    <option value="c">C (GCC)</option>
                    <option value="java">Java 21</option>
                  </select>
                </div>
                <div className="grid grid-cols-2 gap-3">
                  <div className="space-y-1">
                    <label className="text-[10px] font-bold text-white/50 uppercase tracking-widest">
                      Timeout(ms)
                    </label>
                    <input
                      type="number"
                      value={timeLimit}
                      onChange={(e) => setTimeLimit(parseInt(e.target.value))}
                      className="w-full bg-black border border-white/10 px-3 py-2 text-white focus:outline-none focus:border-primary transition-colors font-mono text-sm text-center"
                    />
                  </div>
                  <div className="space-y-1">
                    <label className="text-[10px] font-bold text-white/50 uppercase tracking-widest">
                      RAM(kb)
                    </label>
                    <input
                      type="number"
                      value={memoryLimit}
                      onChange={(e) => setMemoryLimit(parseInt(e.target.value))}
                      className="w-full bg-black border border-white/10 px-3 py-2 text-white focus:outline-none focus:border-primary transition-colors font-mono text-sm text-center"
                    />
                  </div>
                </div>
              </div>

              <div className="space-y-1 pt-2">
                <label className="text-[10px] font-bold text-white/50 uppercase tracking-widest">
                  Payload Source
                </label>
                <textarea
                  value={code}
                  onChange={(e) => setCode(e.target.value)}
                  spellCheck={false}
                  rows={10}
                  className="w-full bg-black border border-white/10 px-4 py-3 text-white/90 focus:outline-none focus:border-primary transition-colors font-mono text-sm resize-y"
                ></textarea>
              </div>

              <div className="pt-4 border-t border-white/10">
                <div className="flex items-center justify-between mb-3">
                  <label className="text-[10px] font-bold text-white/50 uppercase tracking-widest flex items-center gap-2">
                    Test Cases{" "}
                    <span className="px-1 bg-white/10 text-white leading-tight">
                      {testCases.length}
                    </span>
                  </label>
                  <button
                    onClick={() =>
                      setTestCases([
                        ...testCases,
                        {
                          id: testCases.length
                            ? Math.max(...testCases.map((t) => t.id)) + 1
                            : 1,
                          input: "",
                          expectedOutput: "",
                        },
                      ])
                    }
                    className="text-[10px] font-black uppercase tracking-widest bg-white/5 hover:bg-white/10 text-white px-3 py-1.5 transition-colors border border-white/10"
                  >
                    + Add Case
                  </button>
                </div>

                <div className="space-y-4 max-h-[300px] overflow-y-auto pr-2 custom-scrollbar">
                  {testCases.map((tc, idx) => (
                    <div
                      key={tc.id}
                      className="p-4 bg-black border border-white/10 flex flex-col gap-3"
                    >
                      <div className="flex justify-between items-center border-b border-white/5 pb-2">
                        <span className="text-[10px] uppercase font-black text-white/50 tracking-widest">
                          Case #{tc.id}
                        </span>
                        <button
                          onClick={() =>
                            setTestCases(
                              testCases.filter((t) => t.id !== tc.id),
                            )
                          }
                          className="text-red-500 hover:text-red-400 text-[10px] font-black uppercase tracking-widest transition-colors"
                        >
                          Remove
                        </button>
                      </div>
                      <div className="grid grid-cols-2 gap-4">
                        <div className="space-y-1 min-w-0">
                          <label className="text-[10px] text-white/40 uppercase tracking-widest block">
                            Input
                          </label>
                          <textarea
                            value={tc.input}
                            onChange={(e) => {
                              const newCases = [...testCases];
                              newCases[idx].input = e.target.value;
                              setTestCases(newCases);
                            }}
                            spellCheck={false}
                            className="w-full bg-[#111] border border-white/10 px-3 py-1.5 text-white/80 text-xs font-mono focus:outline-none focus:border-primary transition-colors min-h-[40px] resize-y"
                          ></textarea>
                        </div>
                        <div className="space-y-1 min-w-0">
                          <label className="text-[10px] text-white/40 uppercase tracking-widest block">
                            Expected Output
                          </label>
                          <textarea
                            value={tc.expectedOutput}
                            onChange={(e) => {
                              const newCases = [...testCases];
                              newCases[idx].expectedOutput = e.target.value;
                              setTestCases(newCases);
                            }}
                            spellCheck={false}
                            className="w-full bg-[#111] border border-white/10 px-3 py-1.5 text-white/80 text-xs font-mono focus:outline-none focus:border-primary transition-colors min-h-[40px] resize-y"
                          ></textarea>
                        </div>
                      </div>
                    </div>
                  ))}
                  {testCases.length === 0 && (
                    <div className="text-[10px] font-black tracking-widest uppercase text-white/20 text-center py-8 bg-black border border-white/5">
                      No cases added
                    </div>
                  )}
                </div>
              </div>

              <button
                onClick={handleSubmit}
                disabled={isSubmitting}
                className="w-full mt-2 py-4 bg-primary hover:bg-primary-hover text-black font-black text-sm uppercase tracking-widest transition-colors disabled:opacity-50 disabled:cursor-not-allowed flex items-center justify-center gap-2"
              >
                {isSubmitting ? "Dispatching..." : "Dispatch"}
              </button>
            </div>
          </div>
        </div>

        {/* Response / Logs Column */}
        <div className="lg:col-span-6 flex flex-col h-full">
          <div className="bg-surface border border-white/10 flex-1 flex flex-col min-h-[600px] max-h-[85vh]">
            {/* Terminal Header */}
            <div className="flex items-center justify-between px-5 py-4 border-b border-white/10 bg-black">
              <h2 className="text-sm font-black text-white uppercase tracking-widest">
                Execution Trace
              </h2>
              {submissionId && (
                <button
                  onClick={handleCheckStatus}
                  disabled={isChecking}
                  className="px-3 py-1 bg-white/10 hover:bg-white/20 text-white font-black text-[10px] uppercase tracking-widest transition-colors disabled:opacity-50 flex items-center gap-2"
                >
                  {isChecking ? "Polling..." : "Poll Status"}
                </button>
              )}
            </div>

            {/* Submission ID Banner */}
            {submissionId && (
              <div className="bg-primary/10 border-b border-primary/20 px-5 py-3 flex items-center justify-between">
                <span className="text-[10px] text-primary font-black uppercase tracking-widest">
                  Active Trace ID
                </span>
                <span className="font-mono text-[10px] text-primary select-all">
                  {submissionId}
                </span>
              </div>
            )}

            {/* Terminal Body */}
            <div
              ref={logsContainerRef}
              className="flex-1 p-5 overflow-y-auto bg-black font-mono text-xs flex flex-col custom-scrollbar pb-6 relative"
            >
              {statusLog.length === 0 ? (
                <div className="text-white/20 h-full flex items-center justify-center italic text-sm">
                  [ Awaiting payload dispatch... ]
                </div>
              ) : (
                <div className="space-y-4">
                  {statusLog.map((log, i) => (
                    <div
                      key={i}
                      className="flex flex-col gap-1 border-l-2 pl-3 py-0.5"
                      style={{
                        borderLeftColor:
                          log.type === "error"
                            ? "#ef4444"
                            : log.type === "success"
                              ? "#10b981"
                              : "#3b82f6",
                      }}
                    >
                      <div className="flex flex-wrap items-baseline gap-2">
                        <span className="text-white/40 text-[10px]">
                          {log.timestamp.toLocaleTimeString(undefined, {
                            hour12: false,
                            hour: "2-digit",
                            minute: "2-digit",
                            second: "2-digit",
                          })}
                        </span>
                        <span
                          className={`text-[10px] font-black uppercase ${
                            log.type === "error"
                              ? "text-red-400"
                              : log.type === "success"
                                ? "text-green-400"
                                : "text-blue-400"
                          }`}
                        >
                          {log.type}
                        </span>
                        <span className="text-white/80">{log.message}</span>
                      </div>
                      {log.data && (
                        <div className="mt-2 bg-[#111] p-3 border border-white/5 overflow-x-auto custom-scrollbar">
                          <pre className="text-white/60 text-[11px] leading-relaxed">
                            {JSON.stringify(log.data, null, 2)}
                          </pre>
                        </div>
                      )}
                    </div>
                  ))}
                </div>
              )}
            </div>

            {/* Terminal Footer */}
            {statusLog.length > 0 && (
              <div className="p-3 border-t border-white/10 bg-surface flex justify-end">
                <button
                  onClick={() => {
                    setStatusLog([]);
                    setSubmissionId(null);
                  }}
                  className="px-3 py-1.5 text-[10px] uppercase font-black tracking-widest text-white/50 hover:text-white transition-colors border border-transparent hover:border-white/10"
                >
                  Clear Logs
                </button>
              </div>
            )}
          </div>
        </div>
      </div>

      <style jsx global>{`
        .custom-scrollbar::-webkit-scrollbar {
          width: 6px;
          height: 6px;
        }
        .custom-scrollbar::-webkit-scrollbar-track {
          background: #000;
        }
        .custom-scrollbar::-webkit-scrollbar-thumb {
          background: #333;
        }
        .custom-scrollbar::-webkit-scrollbar-thumb:hover {
          background: #555;
        }
      `}</style>
    </div>
  );
}
