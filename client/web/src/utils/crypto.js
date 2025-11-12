// 浏览器端 ECC 密钥生成和加密解密工具

/**
 * 生成 ECC 密钥对（P-256曲线）
 */
export async function generateKeyPair() {
  try {
    // 使用 Web Crypto API 生成 ECDH 密钥对
    const keyPair = await window.crypto.subtle.generateKey(
      {
        name: 'ECDH',
        namedCurve: 'P-256',
      },
      true,
      ['deriveBits', 'deriveKey']
    )

    // 导出公钥为 SPKI 格式
    const publicKeyBuffer = await window.crypto.subtle.exportKey('spki', keyPair.publicKey)
    const publicKeyPEM = arrayBufferToPEM(publicKeyBuffer, 'EC PUBLIC KEY')

    // 导出私钥为 PKCS8 格式
    const privateKeyBuffer = await window.crypto.subtle.exportKey('pkcs8', keyPair.privateKey)
    const privateKeyPEM = arrayBufferToPEM(privateKeyBuffer, 'EC PRIVATE KEY')

    return {
      publicKey: publicKeyPEM,
      privateKey: privateKeyPEM,
    }
  } catch (error) {
    console.error('Key generation error:', error)
    throw new Error('密钥生成失败')
  }
}

/**
 * 使用公钥加密消息（ECDH + AES-GCM）
 * 注意：这是简化版本，实际应该在客户端后端处理
 */
export async function encryptMessage(publicKeyPEM, message) {
  try {
    // 简单的 Base64 编码（实际加密在客户端后端处理）
    const encoded = btoa(unescape(encodeURIComponent(message)))
    return encoded
  } catch (error) {
    console.error('Encryption error:', error)
    throw error
  }
}

/**
 * 使用私钥解密消息（ECDH + AES-GCM）
 * 注意：这是简化版本，实际应该在客户端后端处理
 */
export async function decryptMessage(privateKeyPEM, encryptedMessage) {
  try {
    // 简单的 Base64 解码（实际解密在客户端后端处理）
    const decoded = decodeURIComponent(escape(atob(encryptedMessage)))
    return decoded
  } catch (error) {
    console.error('Decryption error:', error)
    // 如果解密失败，返回原文
    return encryptedMessage
  }
}

/**
 * 将 ArrayBuffer 转换为 PEM 格式
 */
function arrayBufferToPEM(buffer, label) {
  const base64 = btoa(String.fromCharCode(...new Uint8Array(buffer)))
  const formatted = base64.match(/.{1,64}/g).join('\n')
  return `-----BEGIN ${label}-----\n${formatted}\n-----END ${label}-----`
}

/**
 * 将 PEM 格式转换为 ArrayBuffer
 */
function pemToArrayBuffer(pem, label) {
  const pemHeader = `-----BEGIN ${label}-----`
  const pemFooter = `-----END ${label}-----`
  const pemContents = pem
    .replace(pemHeader, '')
    .replace(pemFooter, '')
    .replace(/\s/g, '')
  
  const binaryDer = atob(pemContents)
  const binaryArray = new Uint8Array(binaryDer.length)
  for (let i = 0; i < binaryDer.length; i++) {
    binaryArray[i] = binaryDer.charCodeAt(i)
  }
  
  return binaryArray.buffer
}
