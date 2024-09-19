const crypto = require('crypto');

//block structure 
class Block {
    constructor(index, data, prevHash, nonce){
        this.index = index
        this.data = data
        this.timestamp = new Date().toString()
        this.prevHash = prevHash
        this.nonce = nonce
        this.hash = this.calculateHash()
        
    }
    calculateHash(){
        const res = `${this.index}${this.data}${this.timestamp}${this.prevHash}${this.nonce}`
        const hash = crypto.createHash('sha256')
        hash.update(res)
        return hash.digest('hex')
    }
}

//blockchain structure
class Blockchain{
    constructor(){
        this.blocks = [this.creatingGenesisBlock()]
    }
    //creating the genesis block
    creatingGenesisBlock(){
        return new Block(0,"Genesis Block","",0)
    }

    addBlock(data){
        const prevBlock = this.blocks[this.blocks.length-1]
        let newBlock = new Block(prevBlock.index + 1,data,prevBlock.hash,0)

        //the POW proof of work
        while(!this.isValidHash(newBlock.hash)){
            newBlock.nonce++
            newBlock.hash = newBlock.calculateHash()
        }
        this.blocks.push(newBlock)
       
    }

    //if the pow difficulty is of (eg starts with 00000)
    isValidHash(hash){
        return hash.startsWith('00000')
    }
}

//function main 
function main(){
    const blockchain = new Blockchain()

    blockchain.addBlock("second Block")
    blockchain.addBlock("third block")

    blockchain.blocks.forEach(block => {
        console.log(`Index: ${block.index}`)
        console.log(`Time: ${block.timestamp}`)
        console.log(`Data: ${block.data}`)
        console.log(`Previous hash: ${block.prevHash}`)
        console.log(`Hash: ${block.hash}`)
        console.log(`Nonce: ${block.nonce}\n`)
    })
}
main()