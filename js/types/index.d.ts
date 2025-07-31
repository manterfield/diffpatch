/**
 * Represents a diff operation.
 * Format: [position, deleteCount, additions]
 * - position: The index in the array where the operation starts
 * - deleteCount: Number of elements to delete
 * - additions: Array of elements to add at the position
 */
export type DiffOperation<T = any> = { i: number; d: number; a: T[] };

/**
 * Computes the difference between two arrays and returns a set of operations
 * that can transform the old array into the new array.
 *
 * @param oldArr - The original array
 * @param newArr - The target array to transform to
 * @returns An array of diff operations
 */
export function diff<T>(oldArr: T[], newArr: T[]): DiffOperation<T>[];

/**
 * Applies a set of diff operations to an array to produce a new array.
 *
 * @param arr - The array to apply operations to
 * @param ops - The diff operations to apply
 * @returns A new array with the operations applied
 */
export function applyPatch<T>(arr: T[], ops: DiffOperation<T>[]): T[];
