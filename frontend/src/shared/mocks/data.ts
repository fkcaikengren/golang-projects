import type {
  AuthPayload,
  ProblemDetail,
  ProblemSetSummary,
  ProblemSummary,
  SubmissionItem,
} from '../types/api'

export const mockAuth: AuthPayload = {
  token: 'mock-token',
  user: {
    id: 1,
    email: 'demo@go-oj.dev',
    nickname: 'DemoUser',
    status: 'active',
  },
}

export const mockProblemSets: ProblemSetSummary[] = [
  {
    id: 1,
    name: 'Hot 100',
    slug: 'hot-100',
    description: '覆盖最常见算法题型，适合系统性刷题起步。',
    problemCount: 100,
  },
  {
    id: 2,
    name: '面试高频',
    slug: 'interview-top',
    description: '按面试语境挑选的高频题集合。',
    problemCount: 48,
  },
]

export const mockProblems: ProblemSummary[] = [
  {
    id: 1,
    slug: 'two-sum',
    title: '两数之和',
    difficulty: 'easy',
    tags: ['数组', '哈希表'],
    source: 'LeetCode',
    status: 'solved',
  },
  {
    id: 2,
    slug: 'longest-substring-without-repeating-characters',
    title: '无重复字符的最长子串',
    difficulty: 'medium',
    tags: ['字符串', '滑动窗口'],
    source: 'LeetCode',
    status: 'attempted',
  },
  {
    id: 3,
    slug: 'merge-k-sorted-lists',
    title: '合并 K 个升序链表',
    difficulty: 'hard',
    tags: ['链表', '堆'],
    source: 'LeetCode',
    status: 'unsolved',
  },
]

export const mockProblemDetails: Record<string, ProblemDetail> = {
  'two-sum': {
    ...mockProblems[0],
    description: '给定一个整数数组 nums 和一个整数目标值 target，请你在该数组中找出和为目标值的那两个整数。',
    inputDescription: '第一行输入数组，第二行输入目标值。',
    outputDescription: '输出两个下标，顺序不限。',
    sampleInput: 'nums = [2,7,11,15], target = 9',
    sampleOutput: '[0,1]',
    hint: '可以用哈希表把补数查找降到 O(1)。',
  },
  'longest-substring-without-repeating-characters': {
    ...mockProblems[1],
    description: '给定一个字符串 s ，请你找出其中不含有重复字符的最长子串的长度。',
    inputDescription: '输入一个字符串 s。',
    outputDescription: '输出最长无重复子串长度。',
    sampleInput: 'abcabcbb',
    sampleOutput: '3',
    hint: '窗口右移时同步收缩左边界。',
  },
  'merge-k-sorted-lists': {
    ...mockProblems[2],
    description: '给你一个链表数组，每个链表都已经按升序排列，请将所有链表合并到一个升序链表中。',
    inputDescription: '输入多个升序链表。',
    outputDescription: '输出合并后的升序链表。',
    sampleInput: '[[1,4,5],[1,3,4],[2,6]]',
    sampleOutput: '[1,1,2,3,4,4,5,6]',
    hint: '优先队列可以稳定取当前最小节点。',
  },
}

export const mockSubmissions: SubmissionItem[] = [
  {
    id: 101,
    problemTitle: '两数之和',
    problemSlug: 'two-sum',
    language: 'Go',
    result: 'Accepted',
    submittedAt: '2026-03-27 12:05',
  },
  {
    id: 102,
    problemTitle: '无重复字符的最长子串',
    problemSlug: 'longest-substring-without-repeating-characters',
    language: 'Go',
    result: 'Wrong Answer',
    submittedAt: '2026-03-27 12:18',
  },
]
