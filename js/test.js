import { diff, applyPatch } from "./diff.js";
import fs from "fs";
import path from "path";
import { fileURLToPath } from "url";

const __filename = fileURLToPath(import.meta.url);
const __dirname = path.dirname(__filename);

// Load test file pairs automatically
function loadTestCases() {
  const testDir = path.join(__dirname, "..", "test_documents");
  const files = fs.readdirSync(testDir);
  const cases = [];

  // Find all test pairs
  const pattern = /^test(\d+)_old\.md$/;

  files.forEach((file) => {
    const match = file.match(pattern);
    if (match) {
      const testNum = match[1];
      const oldFile = path.join(testDir, file);
      const newFile = path.join(testDir, `test${testNum}_new.md`);

      // Check if corresponding new file exists
      if (fs.existsSync(newFile)) {
        // Add test cases for each split type
        cases.push({
          name: `test${testNum}_lines`,
          oldFile: oldFile,
          newFile: newFile,
          splitBy: "\n",
        });
        cases.push({
          name: `test${testNum}_sentences`,
          oldFile: oldFile,
          newFile: newFile,
          splitBy: ".",
        });
        cases.push({
          name: `test${testNum}_words`,
          oldFile: oldFile,
          newFile: newFile,
          splitBy: " ",
        });
      }
    }
  });

  return cases;
}

// Load file and split by delimiter
function loadFileSplit(filename, delimiter) {
  const content = fs.readFileSync(filename, "utf8");
  const parts = content.split(delimiter);

  // Filter out empty parts
  return parts.filter((part) => part.trim() !== "");
}

// Simple test runner
function runTests() {
  const cases = loadTestCases();
  let passed = 0;
  let failed = 0;

  console.log(`Running ${cases.length} tests...\n`);

  cases.forEach((testCase) => {
    try {
      // Load old and new arrays
      const oldArr = loadFileSplit(testCase.oldFile, testCase.splitBy);
      const newArr = loadFileSplit(testCase.newFile, testCase.splitBy);

      // Generate patch
      const patch = diff(oldArr, newArr);

      // Apply patch
      const result = applyPatch(oldArr, patch);

      // Verify result equals newArr
      const isEqual = JSON.stringify(result) === JSON.stringify(newArr);

      if (isEqual) {
        console.log(`✓ ${testCase.name}`);
        passed++;
      } else {
        console.log(`✗ ${testCase.name}`);
        console.log(
          `  Expected length: ${newArr.length}, Got length: ${result.length}`
        );
        failed++;
      }
    } catch (error) {
      console.log(`✗ ${testCase.name} - Error: ${error.message}`);
      failed++;
    }
  });

  console.log(`\nResults: ${passed} passed, ${failed} failed`);

  if (failed > 0) {
    process.exit(1);
  }
}

// Benchmark runner
function runBenchmarks() {
  const cases = loadTestCases();

  console.log("Running benchmarks...\n");

  cases.forEach((testCase) => {
    try {
      // Load test data
      const oldArr = loadFileSplit(testCase.oldFile, testCase.splitBy);
      const newArr = loadFileSplit(testCase.newFile, testCase.splitBy);

      // Warm up
      for (let i = 0; i < 10; i++) {
        const patch = diff(oldArr, newArr);
        applyPatch(oldArr, patch);
      }

      // Benchmark diff
      const iterations = 100;
      const startDiff = performance.now();
      for (let i = 0; i < iterations; i++) {
        diff(oldArr, newArr);
      }
      const diffTime = performance.now() - startDiff;

      // Benchmark apply
      const patch = diff(oldArr, newArr);
      const startApply = performance.now();
      for (let i = 0; i < iterations; i++) {
        applyPatch(oldArr, patch);
      }
      const applyTime = performance.now() - startApply;

      // Benchmark combined
      const startCombined = performance.now();
      for (let i = 0; i < iterations; i++) {
        const p = diff(oldArr, newArr);
        applyPatch(oldArr, p);
      }
      const combinedTime = performance.now() - startCombined;

      console.log(`${testCase.name}:`);
      console.log(`  Diff: ${(diffTime / iterations).toFixed(3)}ms/op`);
      console.log(`  Apply: ${(applyTime / iterations).toFixed(3)}ms/op`);
      console.log(`  Combined: ${(combinedTime / iterations).toFixed(3)}ms/op`);
      console.log(`  Patch size: ${JSON.stringify(patch).length} bytes`);
      console.log();
    } catch (error) {
      console.log(`Error benchmarking ${testCase.name}: ${error.message}`);
    }
  });
}

// Check command line arguments
const args = process.argv.slice(2);
if (args.includes("--bench")) {
  runBenchmarks();
} else {
  runTests();
}
