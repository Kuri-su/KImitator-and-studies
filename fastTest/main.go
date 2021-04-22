package main

import (
	"fmt"

	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
)

var b = "$$\n\\left[ \\begin{array}{a} a^l_1 \\\\ ⋮ \\\\ a^l_{d_l} \\end{array}\\right]\n= \\sigma(\n \\left[ \\begin{matrix}\n \tw^l_{1,1} & ⋯  & w^l_{1,d_{l-1}} \\\\\n \t⋮ & ⋱  & ⋮  \\\\\n \tw^l_{d_l,1} & ⋯  & w^l_{d_l,d_{l-1}} \\\\\n \\end{matrix}\\right]  ·\n \\left[ \\begin{array}{x} a^{l-1}_1 \\\\ ⋮ \\\\ ⋮ \\\\ a^{l-1}_{d_{l-1}} \\end{array}\\right] +\n \\left[ \\begin{array}{b} b^l_1 \\\\ ⋮ \\\\ b^l_{d_l} \\end{array}\\right])\n $$"

func main() {
	extensions := parser.CommonExtensions | parser.AutoHeadingIDs
	p := parser.NewWithExtensions(extensions)

	htmlFlags := html.CommonFlags | html.HrefTargetBlank
	opts := html.RendererOptions{Flags: htmlFlags}
	r := html.NewRenderer(opts)

	md := []byte(b)
	output := markdown.ToHTML(md, p, r)

	fmt.Println(string(output))
}

var x = `# [zh-CN] Cryptonote writepaper  

[TOC]

> 目前仅仅包括部分内容, 有翻译不准确或者错误的地方欢迎到 [这里提 Issue 指正](https://github.com/Kuri-su/KBlog)

## 4.2 Definitions (定义)

### 4.2.1 Elliptic curve parameters (椭圆曲线参数)

> As our base signature algorithm we chose to use the fast scheme EdDSA, which is developed and implemented by D.J. Bernstein et al. [18]. Like Bitcoin’s ECDSA it is based on the elliptic curve discrete logarithm problem, so our scheme could also be applied to Bitcoin in future.

作为基础签名算法, 我们选择使用 EdDSA(Edwards-curve Digital Signature Algorithm 爱德华兹曲线数字签名算法)  和 Bitcoin 的 ECDSA (椭圆曲线数字签名算法)一样, EdDSA 也是基于 椭圆曲线离散对数问题, 所以未来 比特币也可以使用 Cryptonote 方案 : )

> Common parameters are:

公共参数如下: 


### 4.2.2 Terminology (术语)

> Enhanced privacy requires a new terminology which should not be confused with Bitcoin entities.

下面会定义一些新的属于, 从而避免和 Bitcoin 的用词相混淆.

> **private ec-key**  is a standard elliptic curve private key: a number $a \in [1,l-1] $
> 
> **public ec-key** is a standard elliptic curve public key: a point $A = aG$;
> 
> **one-time key-air** is a pair of private and public ec-keys;
> 
> **private user key** is a pair (a,b) of two different private ec-keys;
> 
> **tracking key** is a pair (a,B) of private and public ec-key ( where $B=bG $ and $a  b$ // TODO)
> 
> **public user key** is a pair (A,B) of two public ec-keys derived from (a,b);
> 
> **standard address** is a representation of a public user key given into human friendly string with error correction;
> 
> **trucated address** is a representation of the second half (point B) of a public user key given into human friendly string with error correction.

**私有 ec 密钥 (private ec-key)**  是一个 标准的 椭圆曲线 私钥 : 数字 $a \in [1,l-1] $

**公开 ec 密钥 (public ec-key)** 是一个 标准椭圆曲线公钥: 一个点 $A = aG$

**一次性密钥对 (one-time keypair)** 是 一个 基于 椭圆曲线 的 公钥和私钥对


**追踪 密钥 (tracking key)** 是 由 a 的私钥 和 b 公钥组成的密钥对 (a,B),当 $B = bG $  并且 $a \neq b$

**用户的公钥 (public user key)** 是 由 用户的 私钥对(a,b) 派生出的 公钥对 (A,B)


> The transaction structure remains similar to the structure in Bitcoin: every user can choose several independent incoming payments (transactions outputs), sign them with the corresponding private keys and send them to different destinations.

交易结构和 Bitcoin 类似, 每个用户都可以将几个独立的 UTXO 作为输入, 使用对应私钥进行签名, 然后发送到不同的地址上.

> Contrary to Bitcoin’s model, where a user possesses unique private and public key, in the proposed model a sender generates a one-time public key based on the recipient’s address and some random data. In this sense, an incoming transaction for the same recipient is sent to a one-time public key (not directly to a unique address) and only the recipient can recover the corresponding private part to redeem his funds (using his unique private key). The recipient can spend the funds using a ring signature, keeping his ownership and actual spending anonymous. The details of the protocol are explained in the next subsections.


### 4.3 Unlinkable payments (无法连接的交易们)

> Classic Bitcoin addresses, once being published, become unambiguous identifier for incoming payments, linking them together and tying to the recipient’s pseudonyms. If someone wants to receive an “untied” transaction, he should convey his address to the sender by a private channel. If he wants to receive different transactions which cannot be proven to belong to the same owner he should generate all the different addresses and never publish them in his own pseudonym.

在经典的 Bitcoin 网络中, 收款地址一旦公布, 那么就会成为一个明确的标识, 会将这个收款地址和使用者联系起来, 并且计算出通过这个收款地址 转入和转出的资产列表. 如果有人想接受匿名的交易, 那么他应该通过私人渠道将 自己的地址 告诉发件人. 如果他想接受 无法被证明属于同一个所有者的多个交易, 那么他需要生成不同的 Bitcoin 收款地址, 或者永远不要公布这个私人收款地址, 以防止 被和这个私人收款地址关联起来.

![image-20210222131308546](/assets/cryptonote_fig_2.png)

> We propose a solution allowing a user to publish a single address and receive unconditional unlinkable payments. The destination of each CryptoNote output (by default) is a public key, derived from recipient’s address and sender’s random data. The main advantage against Bitcoin is that every destination key is unique by default (unless the sender uses the same data for each of his transactions to the same recipient). Hence, there is no such issue as “address reuse” by design and no observer can determine if any transactions were sent to a specific address or link two addresses together.

我们提出了一个解决方案, 允许用户发布一个单一的收款地址, 也可以做到 接受无法被追踪的转账.


![image-20210222132309803](/assets/cryptonote_fig_3.png)

> First, the sender performs a Diffie-Hellman exchange to get a shared secret from his data and half of the recipient’s address. Then he computes a one-time destination key, using the shared secret and the second half of the address. Two different ec-keys are required from the recipient for these two steps, so a standard CryptoNote address is nearly twice as large as a Bitcoin wallet address. The receiver also performs a Diffie-Hellman exchange to recover the corresponding secret key.


> A standard transaction sequence goes as follows:

一个标准的交易流程如下:

> 1. Alice wants to send a payment to Bob, who has published his standard address. She unpacks the address and gets Bob’s public key (A, B).
> 2. Alice generates a random $r \in [1, l−1]$ and computes a one-time public key $P=H_s(rA)G+B$.
> 3. Alice uses $P$ as a destination key for the output and also packs value $R = rG$ (as a part of the Diffie-Hellman exchange) somewhere into the transaction. Note that she can create other outputs with unique public keys: different recipients’ keys ($A_i , B_i $) imply different $P_i$ even with the same $r$.
>
>![image-20210222132432539](/assets/cryptonote_fig_4.png)

> 4. Alice sends the transaction.
> 2. Bob checks every passing transaction with his private key $(a, b)$, and computes $P' = H_s(aR)G+B$. If Alice’s transaction for with Bob as the recipient was among them, then $aR = arG = rA$ and $P' = P$ .
> 3. Bob can recover the corresponding one-time private key: $x = H_s (aR) + b$, so as $P = xG$. He can spend this output at any time by signing a transaction with $x$.
>
>![image-20210222132514654](/assets/cryptonote_fig_5.png)

4. Alice 发送了一笔交易 到一次性密钥
5. Bob 用它的私钥 $(a,b)$ 检查每笔通过的交易, 并计算  $P' = H_s(aR)G+B$. 如果检查到 发送给 Bob 的 交易, 并且 $aR = arG = rA$ and $P' = P$ .
6. Bob 可以使用他的用户私钥, 来恢复相应的一次性密钥 : $x = H_s (aR) + b$, 得到 $P=xG$ . 他可以在任何时候用 $x$ 签署一个 交易来花费这笔 UTXO.

> As a result Bob gets incoming payments, associated with one-time public keys which are unlinkable for a spectator. Some additional notes:

结果 Bob 得到了 这笔转入资金. 由于这个一次性公钥仅仅使用一次, 所以这个这笔交易对 旁观者来说是无法连接的, 无法推断出这笔交易的双方. 除了上面信息之外, 这里还有一些附加说明

> * When Bob “recognizes” his transactions (see step 5) he practically uses only half of his private information: $(a, B)$. This pair, also known as the tracking key, can be passed to a third party (Carol). Bob can delegate her the processing of new transactions. Bob doesn’t need to explicitly trust Carol, because she can’t recover the one-time secret key $p$ without Bob’s full private key $(a, b)$. This approach is useful when Bob lacks bandwidth or computation power (smartphones, hardware wallets etc.).
>* In case Alice wants to prove she sent a transaction to Bob’s address she can either disclose $r$ or use any kind of zero-knowledge protocol to prove she knows $r$ (for example by signing the transaction with $r$).
> * If Bob wants to have an audit compatible address where all incoming transaction are linkable, he can either publish his tracking key or use a truncated address. That address represent only one public ec-key $B$, and the remaining part required by the protocol is derived from it as follows: $a=H_s(B)$ and $A=H_s(B)G$. In both cases every person is able to “recognize” all of Bob’s incoming transaction, but, of course, none can spend the funds enclosed within them without the secret key b.

* 如果 Bob 想拥有一个 可以被审计的地址, 所有的 传入交易都是可以链接的. 那么他可以发布 他的 跟踪密钥 $(a,B)$ ,或者 一个截断的地址 B, 协议所需的其他部分由 B 派生,  $a=H_s(B)$ and $A=H_s(B)G$. 但是, 由于没有密钥 b, 所以谁也不能花掉其中的资金.

### 4.4 One-time ring signatures (一次性环签名)

> A protocol based on one-time ring signatures allows users to achieve unconditional unlinkability. Unfortunately, ordinary types of cryptographic signatures permit to trace transactions to their respective senders and receivers. Our solution to this deficiency lies in using a different signature type than those currently used in electronic cash systems.


> We will first provide a general description of our algorithm with no explicit reference to electronic cash.

首先我们先对 算法进行一般性描述, 而不提及电子现金系统.

> A one-time ring signature contains four algorithms: (**GEN**, **SIG**, **VER**, **LNK**):
>
> * GEN: takes public parameters and outputs an ec-pair $(P, x)$ and a public key $I$.
> * SIG: takes a message m, a set $S'$ of public keys $\{P_i\}_{i\neq s}$ , a pair $(P_s , x_s)$ and outputs a signature $σ$ and a set $S = S' ∪ \{P_s \}$.
> * VER: takes a message m, a set S, a signature σ and outputs “true” or “false”.
> * LNK: takes a set $I = {I i }$, a signature $σ$ and outputs “linked” or “indep”.

一次性环签名包含四种算法:  (**GEN**, **SIG**, **VER**, **LNK**):

* GEN: 接受公共参数 然后 输出 密钥对 $(P,x)$ 和 公钥 $I$
* SIG: 接受如下参数, 输出签名 $σ$ 和 一组 $S = S' ∪ \{P_s \}$
  * 一个信息 m
  * 一组 关于 $S'$ 的公钥  $\{P_i\}_{i\neq s}$
* LNK: 接受一对 $I$=$\{I_i\}$, 一个 签名 σ, 然后输出 “linked” or “indep”

> The idea behind the protocol is fairly simple: a user produces a signature which can be
> checked by a set of public keys rather than a unique public key. The identity of the signer is
> indistinguishable from the other users whose public keys are in the set until the owner produces
> a second signature using the same keypair.
>
> ![image-20210224134609552](/assets/cryptonote_fig_6.png)
>
> * GEN: The signer picks a random secret key x ∈ [1, l − 1] and computes the corresponding
>   public key P = xG. Additionally he computes another public key I = xH p (P ) which we will
>   call the “key image”.
> * SIG: The signer generates a one-time ring signature with a non-interactive zero-knowledge
>   proof using the techniques from [21]. He selects a random subset S 0 of n from the other users’
>   public keys P i , his own keypair (x, P ) and key image I. Let 0 ≤ s ≤ n be signer’s secret index
>   in S (so that his public key is P s ).

> He picks a random {q i | i = 0 . . . n} and {w i | i = 0 . . . n, i 6 = s} from (1 . . . l) and applies the
> following transformations:
> $$
> L_i=\begin{cases}
> q_iG & \text{if } i = s; \\
> q_iG + w_iP_i, & \text{if } i \neq s \\
> \end{cases}
> \\
> R_i=\begin{cases}
> q_iH_p(P_i), & \text{if } i =s \\
> q_iH_p(P_i) + w_iI, & \text{if } i \neq s\\
> \end{cases}
> $$
>
> The next step is getting the non-interactive challenge:
> $$
> c = H_s (m, L _1 , . . . , L_n , R_1 , . . . , R_n )
> $$



Finally the signer computes the response:

$$
c_i=\begin{cases}
w_i && \text{if } i \neq s; \\
c - \sum_{i=0}^n c_i & mod \space l, & \text{if } i = s \\
\end{cases}
\\
r_i=\begin{cases}
q_i && \text{if } i \neq s; \\
q_s - c_sx & mod \space l, & \text{if } i = s \\
\end{cases}
$$

>
>
> The resulting signature is$ σ = (I, c_1 , . . . , c_n , r_1 , . . . , r_n )$.gotjib
>
> VER: The verifier checks the signature by applying the inverse transformations:
> $$
> \begin{cases}
> L_i^\prime = r_iG + c_iP_i \\
> R_i^\prime = r_iH_p(P_i)+c_iI
> \end{cases}
> $$
>
>
>
> Finally, the verifier checks ifnP?c i = H s (m, L 0 0 , . . . , L 0 n , R 0 0 , . . . , R n 0 ) mod l
> If this equality is correct, the verifier runs the algorithm LNK. Otherwise the verifier rejects
> the signature.
> LNK: The verifier checks if I has been used in past signatures (these values are stored in the
> set I). Multiple uses imply that two signatures were produced under the same secret key.
> The meaning of the protocol: by applying L-transformations the signer proves that he knows
> such x that at least one P i = xG. To make this proof non-repeatable we introduce the key image
> as I = xH p (P ). The signer uses the same coefficients (r i , c i ) to prove almost the same statement:
> he knows such x that at least one H p (P i ) = I · x −1 .
> If the mapping x → I is an injection:
>
> 1. Nobody can recover the public key from the key image and identify the signer;
> 2. The signer cannot make two signatures with different I’s and the same x.
> A full security analysis is provided in Appendix A.

### 4.5 Standard CryptoNote transaction (标准 CryptoNote 转账)

> By combining both methods (unlinkable public keys and untraceable ring signature) Bob achieves new level of privacy in comparison with the original Bitcoin scheme. It requires him to store only one private key (a, b) and publish (A, B) to start receiving and sending anonymous transactions.
> While validating each transaction Bob additionally performs only two elliptic curve multi-plications and one addition per output to check if a transaction belongs to him. For his every output Bob recovers a one-time keypair (p i , P i ) and stores it in his wallet. Any inputs can be circumstantially proved to have the same owner only if they appear in a single transaction. In fact this relationship is much harder to establish due to the one-time ring signature.
> With a ring signature Bob can effectively hide every input among somebody else’s; all possible spenders will be equiprobable, even the previous owner (Alice) has no more information than any observer.
> When signing his transaction Bob specifies n foreign outputs with the same amount as his output, mixing all of them without the participation of other users. Bob himself (as well as anybody else) does not know if any of these payments have been spent: an output can be used in thousands of signatures as an ambiguity factor and never as a target of hiding. The double spend check occurs in the LNK phase when checking against the used key images set.
> Bob can choose the ambiguity degree on his own: n = 1 means that the probability he has spent the output is 50% probability, n = 99 gives 1%. The size of the resulting signature increases linearly as O(n + 1), so the improved anonymity costs to Bob extra transaction fees. He also can set n = 0 and make his ring signature to consist of only one element, however this will instantly reveal him as a spender.

`
