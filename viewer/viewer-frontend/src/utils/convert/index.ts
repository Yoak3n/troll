const XOR_CODE = 23442827791579n;
const MASK_CODE = 2251799813685247n;
const MAX_AID = 1n << 51n;
const BASE = 58n;

const data = 'FcwAPNKTMug3GV5Lj7EJnHpWsx4tb8haYeviqBz6rkCy12mUSDQX9RdoZf';

export function av2bv(aid: number) {
  const bytes = ['B', 'V', '1', '0', '0', '0', '0', '0', '0', '0', '0', '0'];
  let bvIndex = bytes.length - 1;
  let tmp = (MAX_AID | BigInt(aid)) ^ XOR_CODE;
  while (tmp > 0) {
    const char = data[Number(tmp % BigInt(BASE))];
    if (char === undefined) throw new Error('Invalid index for data');
    bytes[bvIndex] = char;
    tmp = tmp / BASE;
    bvIndex -= 1;
  }
  const tmp3 = bytes[3];
  const tmp9 = bytes[9];
  if (tmp3 === undefined || tmp9 === undefined) throw new Error('Invalid bytes');
  bytes[3] = tmp9;
  bytes[9] = tmp3;
  const tmp4 = bytes[4];
  const tmp7 = bytes[7];
  if (tmp4 === undefined || tmp7 === undefined) throw new Error('Invalid bytes');
  bytes[4] = tmp7;
  bytes[7] = tmp4;
  return bytes.join('') as `BV1${string}`;
}

export function bv2av(bvid: `BV1${string}`) {
  const bvidArr = Array.from<string>(bvid);
  const tmp3 = bvidArr[3];
  const tmp9 = bvidArr[9];
  if (tmp3 === undefined || tmp9 === undefined) throw new Error('Invalid bvid');
  bvidArr[3] = tmp9;
  bvidArr[9] = tmp3;
  const tmp4 = bvidArr[4];
  const tmp7 = bvidArr[7];
  if (tmp4 === undefined || tmp7 === undefined) throw new Error('Invalid bvid');
  bvidArr[4] = tmp7;
  bvidArr[7] = tmp4;
  bvidArr.splice(0, 3);
  const tmp = bvidArr.reduce((pre, bvidChar) => pre * BASE + BigInt(data.indexOf(bvidChar)), 0n);
  return Number((tmp & MASK_CODE) ^ XOR_CODE);
}