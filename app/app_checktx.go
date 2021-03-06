package app

import (
	"github.com/tendermint/abci/types"
	"encoding/hex"
	"log"
)

/**
 * should be very long validating...
 */
func Exec_CheckTx(tx []byte) types.Result {
	var err error

	// decode hex string to json string
	arr, err := hex.DecodeString(string(tx[:]))
	if (err != nil) {
		log.Println(err.Error())
		return types.ErrEncodingError
	}
	jsonstring:= string(arr[:])

	// decode jsonstring to OperationEnvelope and Operation
	envelope, operation , err := UnMarshal(jsonstring)
	if (err != nil) {
		log.Println(err.Error())
		return types.ErrEncodingError
	}

	///////////////////////////////////
	// Validate the envelope
	///////////

	// decode PubKey
	bPubKey, err := hex.DecodeString(envelope.Pubkey)
	if (err != nil) {
		log.Println(err.Error())
		return types.ErrEncodingError
	}

	pubKey, err := GetPubKeyFromBytes(bPubKey)
	if (err != nil) {
		log.Println(err.Error())
		return types.ErrInternalError
	}

	// decode SignatureEd25519
	bSignature, err := hex.DecodeString(envelope.Signature)
	if (err != nil) {
		log.Println(err.Error())
		return types.ErrEncodingError
	}

	signature, err := GetSignatureFromBytes(bSignature)
	if (err != nil) {
		log.Println(err.Error())
		return types.ErrInternalError
	}

	// verify signature
	verify := pubKey.VerifyBytes([]byte(envelope.Operation), signature)
	if (!verify) {
		log.Println("Can not verify signature")
		return types.ErrUnauthorized
	}

	// address of the user who make the transaction
	address := pubKey.Address()

	// check if Address existing
	if (envelope.Type != "AccountCreateOperation") {
		_, err = TreeGetAccount(GlogGlobal.GlogApp.State, address)
		if (err != nil) {
			log.Println(err.Error())
			return types.ErrInternalError
		}
	}

	// validate the Fee
	if (envelope.Fee < 0) {
		log.Println("Fee value is invalid")
		return types.ErrInternalError
	}

	///////////////////////////////////
	// Validate the operation
	///////////

	switch operation.(type) {
	case AccountCreateOperation:
		_, ok := operation.(AccountCreateOperation)
		if (!ok) {
			log.Println("Can not cast operation to AccountCreateOperation")
			return types.ErrInternalError
		}

		break
	case SendTokenOperation:
		_, ok := operation.(SendTokenOperation)
		if (!ok) {
			log.Println("Can not cast operation to SendTokenOperation")
			return types.ErrInternalError
		}

		break
	case PostCreateOperation:
		_, ok := operation.(PostCreateOperation)
		if (!ok) {
			log.Println("Can not cast operation to PostCreateOperation")
			return types.ErrInternalError
		}

		break
	case PostEditOperation:
		_, ok := operation.(PostEditOperation)
		if (!ok) {
			log.Println("Can not cast operation to PostEditOperation")
			return types.ErrInternalError
		}

		break
	case CommentCreateOperation:
		_, ok := operation.(CommentCreateOperation)
		if (!ok) {
			log.Println("Can not cast operation to CommentCreateOperation")
			return types.ErrInternalError
		}
		break
	default:
		log.Println("Unknow Operation!")
		return types.ErrUnknownRequest
	}

	//TODO: more validation here or the chain get fucked!
	//TODO: make a list of security threats
	//TODO: write tests (jake) (but tuan please check my tests!)


	//var buf bytes.Buffer        // Stand-in for a network connection
	//enc := gob.NewEncoder(&buf) // Will write to network.
	//dec := gob.NewDecoder(&buf) // Will read from network.
	//
	//err = enc.Encode(obj)
	//if (err != nil) {
	//	log.Println(err.Error())
	//	return types.ErrInternalError
	//}


 	//TODO: avoid double spending within current block

	return types.OK
}
