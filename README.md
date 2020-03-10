# Asynchronous Consensus-Free Transaction Systems

## Topology Configuration Spec

```ts
[
    {
        type: "server" | "client"
        address: Blob
        port: number
        publicKey: Blob
        privateKey: Blob
        agentType?: "randomSingle"
    }
]
```
`Blob`s are encoded in a Base64 string.
