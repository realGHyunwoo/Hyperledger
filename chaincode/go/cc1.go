package main

import (
   "bytes"
   "crypto/sha256"
   "encoding/hex"
   "encoding/json"
   "fmt"
   "strconv"
   "strings"
   "time"

   "github.com/hyperledger/fabric/core/chaincode/lib/cid"
   "github.com/hyperledger/fabric/core/chaincode/shim"
   pb "github.com/hyperledger/fabric/protos/peer"
)

type SmartContract struct{}

var walletflag bool = false

// 체인코드 설치 및 업데이트 할 때 실행
func (s *SmartContract) Init(APIstub shim.ChaincodeStubInterface) pb.Response {
   return shim.Success(nil)
}

// 체인코드 호출
func (s *SmartContract) Invoke(APIstub shim.ChaincodeStubInterface) pb.Response {
   function, args := APIstub.GetFunctionAndParameters()

   if function == "initWallet" {
      return s.initWallet(APIstub)
   } else if function == "setWallet" {
      return s.setWallet(APIstub, args)
   } else if function == "getWallet" {
      return s.getWallet(APIstub, args)
   } else if function == "getCertificate" {
      return s.getCertificate(APIstub, args)
   } else if function == "setCertificate" {
      return s.setCertificate(APIstub, args)
   } else if function == "updateCertificate" {
      return s.updateCertificate(APIstub, args)
   } else if function == "returnCertificate" {
      return s.returnCertificate(APIstub, args)
   } else if function == "updateRating" {
      return s.updateRating(APIstub, args)
   } else if function == "setJobPosting" {
      return s.setJobPosting(APIstub, args)
   } else if function == "verifyConditions" {
      return s.verifyConditions(APIstub, args)
   } else if function == "setFreelancer" {
      return s.setFreelancer(APIstub, args)
   } else if function == "setRating" {
      return s.setRating(APIstub, args)
   } else if function == "getJobPosting" {
      return s.getJobPosting(APIstub)
   } else if function == "deleteJobPosting" {
      return s.deleteJobPosting(APIstub, args)
   } else if function == "setApply" {
      return s.setApply(APIstub, args)
   }

   fmt.Println("Please check your function : " + function)
   return shim.Error("Unknown function")
}

func main() {
   err := shim.Start(new(SmartContract))
   if err != nil {
      fmt.Printf("Error starting Simple chaincode: %s", err)
   }
}

type Wallet struct {
   Name  string `json:"name"`
   ID    string `json:"id"`
   Token string `json:"token"`
}

type Document struct {
   WalletID        string            `json:"walletid"`
   PSWORD          string            `json:"psword"`
   University_code string            `json:"university_code"`
   Certification   map[string]string `json:"certification"`
   Grade           string            `json:"grade"`
   Count           string            `json:"count"`
}

type DocsKey struct {
   Key string
   Idx int
}

type List struct {
   list map[string]string
}

func (s *SmartContract) initWallet(APIstub shim.ChaincodeStubInterface) pb.Response {

   if walletflag {
      return shim.Error("Already initWallet")
   }

   //Declare wallets
   company := Wallet{Name: "Samsung", ID: "Company", Token: "10000000"}
   freelancer := Wallet{Name: "Freelancer", ID: "jjy", Token: "0"}

   // Convert seller to []byte
   CompanysJSONBytes, _ := json.Marshal(company)
   err := APIstub.PutState(company.ID, CompanysJSONBytes)
   if err != nil {
      return shim.Error("Failed to create asset " + company.Name)
   }
   // Convert customer to []byte
   FreelancersJSONBytes, _ := json.Marshal(freelancer)
   err = APIstub.PutState(freelancer.ID, FreelancersJSONBytes)
   if err != nil {
      return shim.Error("Failed to create asset " + freelancer.Name)
   }

   walletflag = true

   return shim.Success(nil)
}

// param - Name, ID, skey
func (s *SmartContract) setWallet(APIstub shim.ChaincodeStubInterface, args []string) pb.Response {

   mspid, _ := cid.GetMSPID(APIstub)

   if mspid != "CertOrg" {
      return shim.Error("You do not have permission")
   }

   if len(args) != 3 {
      return shim.Error("Incorrect number of arguments. Expecting 3")
   }
   var wallet = Wallet{Name: args[0], ID: args[1], Token: args[2]}

   WalletasJSONBytes, _ := json.Marshal(wallet)
   err := APIstub.PutState(wallet.ID, WalletasJSONBytes)
   if err != nil {
      return shim.Error("Failed to create asset " + wallet.Name)
   }

   return shim.Success(nil)
}

// param - wallet key(Wallet.ID)
func (s *SmartContract) getWallet(APIstub shim.ChaincodeStubInterface, args []string) pb.Response {
   if len(args) != 1 {
      return shim.Error("Incorrect number of arguments. Expecting 1")
   }

   walletAsBytes, err := APIstub.GetState(args[0])
   if err != nil {
      fmt.Println(err.Error())
   } else if walletAsBytes == nil {
      return shim.Error("Could not Find ID")
   }

   wallet := Wallet{}
   json.Unmarshal(walletAsBytes, &wallet)

   var buffer bytes.Buffer
   buffer.WriteString("[")
   bArrayMemberAlreadyWritten := false

   if bArrayMemberAlreadyWritten == true {
      buffer.WriteString(",")
   }
   buffer.WriteString("{\"Name\":")
   buffer.WriteString("\"")
   buffer.WriteString(wallet.Name)
   buffer.WriteString("\"")

   buffer.WriteString(", \"ID\":")
   buffer.WriteString("\"")
   buffer.WriteString(wallet.ID)
   buffer.WriteString("\"")

   buffer.WriteString(", \"Token\":")
   buffer.WriteString("\"")
   buffer.WriteString(wallet.Token)
   buffer.WriteString("\"")

   buffer.WriteString("}")
   bArrayMemberAlreadyWritten = true
   buffer.WriteString("]\n")

   return shim.Success(buffer.Bytes())
}

func generateKey(APIstub shim.ChaincodeStubInterface, key string) []byte {

   var isFirst bool = false

   docskeyAsBytes, err := APIstub.GetState(key)
   if err != nil {
      fmt.Println(err.Error())
   }

   docskey := DocsKey{}
   json.Unmarshal(docskeyAsBytes, &docskey)
   var tempIdx string
   tempIdx = strconv.Itoa(docskey.Idx)
   fmt.Println(docskey)
   fmt.Println("Key is " + strconv.Itoa(len(docskey.Key)))
   if len(docskey.Key) == 0 || docskey.Key == "" {
      isFirst = true
      docskey.Key = "freelancerKey"
   }
   if !isFirst {
      docskey.Idx = docskey.Idx + 1
   }

   fmt.Println("Last DocsKey is " + docskey.Key + " : " + tempIdx)

   returnValueBytes, _ := json.Marshal(docskey)

   return returnValueBytes
}

// params - Name, Age, WalletID, U_code Cert
func (s *SmartContract) setCertificate(APIstub shim.ChaincodeStubInterface, args []string) pb.Response {
   mspid, _ := cid.GetMSPID(APIstub)

   if mspid != "CertOrg" {
      return shim.Error("You do not have permission")
   }

   if len(args) != 6 && len(args) != 4 {
      return shim.Error("Incorrect number of arguments. Expecting 6")
   }
   // Freelancer의 문서 Key생성

   walletid, psword, university_code, certification := args[0], args[1], args[2], args[3]

   var docskey = DocsKey{}
   json.Unmarshal(generateKey(APIstub, "latestKey"), &docskey)
   keyidx := strconv.Itoa(docskey.Idx)

   slicestr := strings.Split(certification, ",")

   certificate := make(map[string]string)

   for _, valuse := range slicestr {
      if strings.Contains(valuse, ":") {
         key_valuse := strings.Split(valuse, ":")
         certificate[key_valuse[0]] = key_valuse[1]
         continue
      }
      certificate[valuse] = "ture"
   }

   hash := sha256.New()

   hash.Write([]byte(psword))

   word := hash.Sum(nil)

   psStr := hex.EncodeToString(word)

   var document = Document{}
   if len(args) == 6 {
      document.WalletID = walletid
      document.PSWORD = psStr
      document.University_code = university_code
      document.Certification = certificate
      document.Grade = args[4]
      document.Count = args[5]
   } else if len(args) == 4 {
      document.WalletID = walletid
      document.PSWORD = psStr
      document.University_code = university_code
      document.Certification = certificate
      document.Grade = "0"
      document.Count = "0"
   }
   documentAsJSONBytes, _ := json.Marshal(document)

   var keyString = docskey.Key + keyidx

   err := APIstub.PutState(keyString, documentAsJSONBytes)
   if err != nil {
      return shim.Error(fmt.Sprintf("Failed to record document catch: %s", "Key"))
   }

   docskeyAsBytes, _ := json.Marshal(docskey)
   APIstub.PutState("latestKey", docskeyAsBytes)

   return shim.Success(nil)
}

// params - freelancerdocskey
func (s *SmartContract) getCertificate(APIstub shim.ChaincodeStubInterface, args []string) pb.Response {
   if len(args) != 1 {
      return shim.Error("Incorrect number of arguments. Expecting 1")
   }

   documentKey := args[0]

   // FreelancerID로 Document 불러오기
   documentAsBytes, err := APIstub.GetState(documentKey)
   if err != nil {
      fmt.Println(err.Error())
   } else if documentAsBytes == nil {
      return shim.Error("Could not Find document")
   }

   document := Document{}
   json.Unmarshal(documentAsBytes, &document)

   wokeCount, _ := strconv.ParseFloat(document.Count, 64)
   workgrade, _ := strconv.ParseFloat(document.Grade, 64)

   var grade float64 = workgrade / wokeCount

   var s2 string
   s2 = strconv.FormatFloat(grade, 'f', 1, 64) // f, fmt, prec, bitSize

   var buffer bytes.Buffer
   buffer.WriteString("[")
   bArrayMemberAlreadyWritten := false

   if bArrayMemberAlreadyWritten == true {
      buffer.WriteString(",")
   }

   buffer.WriteString("{\"WalletID\":")
   buffer.WriteString("\"")
   buffer.WriteString(document.WalletID)

   buffer.WriteString(", \"University_code\":")
   buffer.WriteString("\"")
   buffer.WriteString(document.University_code)
   buffer.WriteString("\"")

   buffer.WriteString(", \"Certification\":")
   buffer.WriteString("\"")
   for key, values := range document.Certification {
      buffer.WriteString(key)
      buffer.WriteString(":")
      buffer.WriteString(values)
      buffer.WriteString(",")
   }
   buffer.WriteString("\"")

   buffer.WriteString(", \"Grade\":")
   buffer.WriteString("\"")
   buffer.WriteString(s2)
   buffer.WriteString("\"")

   buffer.WriteString(", \"Count\":")
   buffer.WriteString("\"")
   buffer.WriteString(document.Count)
   buffer.WriteString("\"")

   buffer.WriteString("}")
   bArrayMemberAlreadyWritten = true
   buffer.WriteString("]\n")

   return shim.Success(buffer.Bytes())
}

// params  - freelancerDocskey
func (s *SmartContract) returnCertificate(APIstub shim.ChaincodeStubInterface, args []string) pb.Response {
   mspid, _ := cid.GetMSPID(APIstub)

   if mspid != "CertOrg" {
      return shim.Error("You do not have permission")
   }

   if len(args) != 1 {
      return shim.Error("Incorrect number of arguments. Expecting 1")
   }

   documentKey := args[0]

   // FreelancerID로 Document 불러오기
   documentAsBytes, err := APIstub.GetState(documentKey)
   if err != nil {
      fmt.Println(err.Error())
   } else if documentAsBytes == nil {
      return shim.Error("Could not Find document")
   }

   return shim.Success(documentAsBytes)
}

// params - freelancerDocskey, 자격증 내용
func (s *SmartContract) updateCertificate(APIstub shim.ChaincodeStubInterface, args []string) pb.Response {
   mspid, _ := cid.GetMSPID(APIstub)

   if mspid != "CertOrg" {
      return shim.Error("You do not have permission")
   }

   // 받는 파라미터가 2개가 아닐경우 오류처리
   if len(args) != 2 {
      return shim.Error("Incorrect number of arguments. Expecting 2")
   }

   documentKey, addCertification := args[0], args[1]

   documentbytes, err := APIstub.GetState(documentKey)
   // docs에 조회가 되지 않거나 불러들이는데 실패하면 에러 처리
   if err != nil {
      return shim.Error("Could not locate document")
   } else if documentbytes == nil {
      return shim.Error("Could not Find document")
   }

   // Document구조체를 생성하고 bytes형식의 데이터를 json형식으로 변환한다
   document := Document{}
   json.Unmarshal(documentbytes, &document)

   // Certification 슬라이스에 새로운 자격증 추가
   strtemp := addCertification

   slicestr := strings.Split(strtemp, ",")

   for _, valuse := range slicestr {
      if strings.Contains(valuse, ":") {
         key_valuse := strings.Split(valuse, ":")
         document.Certification[key_valuse[0]] = key_valuse[1]
         continue
      }
      document.Certification[valuse] = "ture"
   }

   documentbytes, _ = json.Marshal(document)
   err2 := APIstub.PutState(documentKey, documentbytes)
   if err2 != nil {
      return shim.Error(fmt.Sprintf("Failed to change document price: %s", documentKey))
   }
   return shim.Success(nil)
}

// param - CompanyID, FreelancerID, grade, Token
func (s *SmartContract) updateRating(APIstub shim.ChaincodeStubInterface, args []string) pb.Response {

   var workerToken int
   var wokeCount int
   var companyToken int
   var requirepay int
   var workgrade int
   var addgrade int
   // check param
   if len(args) != 4 {
      return shim.Error("Incorrect number of parameters 4")
   }

   workerDocskey, companywalletKey, grade, token := args[0], args[1], args[2], args[3]
   // check token(예치금)

   // get Freelancer Document
   documentAsBytes, err := APIstub.GetState(workerDocskey)
   if err != nil {
      fmt.Println(err.Error())
   } else if documentAsBytes == nil {
      return shim.Error("Could not Find document")
   }

   document := Document{}
   json.Unmarshal(documentAsBytes, &document)

   // get Freelancer token
   walletAsBytes, err := APIstub.GetState(document.WalletID)
   if err != nil {
      fmt.Println(err.Error())
   } else if walletAsBytes == nil {
      return shim.Error("Could not Find ID")
   }

   wallet := Wallet{}
   json.Unmarshal(walletAsBytes, &wallet)

   companywalletAsBytes, err := APIstub.GetState(companywalletKey)
   if err != nil {
      fmt.Println(err.Error())
   } else if walletAsBytes == nil {
      return shim.Error("Could not Find ID")
   }

   comwallet := Wallet{}
   json.Unmarshal(companywalletAsBytes, &comwallet)

   // calculate token
   workerToken, _ = strconv.Atoi(wallet.Token)
   requirepay, _ = strconv.Atoi(token)
   wallet.Token = strconv.Itoa(workerToken + requirepay)

   companyToken, _ = strconv.Atoi(comwallet.Token)
   comwallet.Token = strconv.Itoa(companyToken - requirepay)

   // wallet State Update
   walletUpdate, _ := json.Marshal(wallet)
   APIstub.PutState(document.WalletID, walletUpdate)

   comwalletUpdate, _ := json.Marshal(comwallet)
   APIstub.PutState(companywalletKey, comwalletUpdate)

   wokeCount, _ = strconv.Atoi(document.Count)
   wokeCount += 1
   workgrade, _ = strconv.Atoi(document.Grade)
   addgrade, _ = strconv.Atoi(grade)

   workgrade += addgrade

   document.Grade = strconv.Itoa(workgrade)
   document.Count = strconv.Itoa(wokeCount)

   docsUpdate, _ := json.Marshal(document)
   APIstub.PutState(workerDocskey, docsUpdate)

   return shim.Success([]byte("transfer Success"))
}

type Offer struct {
   Filed          string            `json:"filed"`
   C_name         string            `json:"c_name"` // 회사명
   C_PSWORD       string            `json:"c_psword"`
   Requirement    map[string]string `json:"requirement"`    // 자격조건
   Remuneration   string            `json:"remuneration"`   // 보수
   Volunteer      []Document        `json:"volunteer"`      // 지원자 명단
   ContractPeriod string            `json:"contractperiod"` // 계약기간
   State          bool              `json:"state"`
   Timestamp      time.Time         `json:"timestamp"`
}

type OfferKey struct {
   Key string
   Idx int
}

func generateKey2(APIstub shim.ChaincodeStubInterface, key string) []byte {

   var isFirst bool = false

   offerkeyAsBytes, err := APIstub.GetState(key)
   if err != nil {
      fmt.Println(err.Error())
   }

   offerkey := OfferKey{}
   json.Unmarshal(offerkeyAsBytes, &offerkey)
   var tempIdx string
   tempIdx = strconv.Itoa(offerkey.Idx)
   fmt.Println(offerkey)
   fmt.Println("Key is " + strconv.Itoa(len(offerkey.Key)))
   if len(offerkey.Key) == 0 || offerkey.Key == "" {
      isFirst = true
      offerkey.Key = "postKey"
   }
   if !isFirst {
      offerkey.Idx = offerkey.Idx + 1
   }

   fmt.Println("Last OfferKey is " + offerkey.Key + " : " + tempIdx)

   returnValueBytes, _ := json.Marshal(offerkey)

   return returnValueBytes
}

// 구인 공고
// params - 분야 ,CompanyName, 요구조건, 보수, 계약기간
func (s *SmartContract) setJobPosting(APIstub shim.ChaincodeStubInterface, args []string) pb.Response {
   if len(args) != 6 {
      return shim.Error("Incorrect number of arguments. Expecting 6")
   }
   field, companyName, psword, strtemp, token, enddate := args[0], args[1], args[2], args[3], args[4], args[5]

   var offerkey = OfferKey{}
   json.Unmarshal(generateKey2(APIstub, "offerlatesKey"), &offerkey)
   keyidx := strconv.Itoa(offerkey.Idx)
   fmt.Println("Key : " + offerkey.Key + ", Idx : " + keyidx)

   slicestr := strings.Split(strtemp, ",")

   requir := make(map[string]string)

   for _, valuse := range slicestr {
      if strings.Contains(valuse, ":") {
         key_valuse := strings.Split(valuse, ":")
         requir[key_valuse[0]] = key_valuse[1]
         continue
      }
      requir[valuse] = "ture"
   }

   var volunteer []Document
   date := time.Now()

   hash := sha256.New()

   hash.Write([]byte(psword))

   word := hash.Sum(nil)

   psStr := hex.EncodeToString(word)

   var offer = Offer{Filed: field, C_name: companyName, C_PSWORD: psStr, Requirement: requir, Remuneration: token, Volunteer: volunteer, ContractPeriod: enddate, State: true, Timestamp: date}
   offerAsJSONBytes, _ := json.Marshal(offer)

   var keyString = offerkey.Key + keyidx
   fmt.Println("offerkey is" + keyString)

   err := APIstub.PutState(keyString, offerAsJSONBytes)
   if err != nil {
      return shim.Error(fmt.Sprintf("Failed to record Offer catch : %s", offerkey.Key))
   }

   offerkeyAsBytes, _ := json.Marshal(offerkey)
   APIstub.PutState("offerlatesKey", offerkeyAsBytes)

   return shim.Success(nil)
}

// 고용
func (s *SmartContract) setFreelancer(APIstub shim.ChaincodeStubInterface, args []string) pb.Response {
   if len(args) != 1 {
      return shim.Error("Incorrect number of arguments. Expecting 1")
   }

   offerkey := args[0]

   offer := Offer{}
   offerBytes, _ := APIstub.GetState(offerkey)
   json.Unmarshal(offerBytes, &offer)

   if !offer.State {
      return shim.Error("Already setFreelancerd worker")
   }

   if len(offer.Volunteer) == 0 {
      return shim.Error("zero")
   }

   if time.Since(offer.Timestamp) < 36000 {
      return shim.Error("today not setFreelancer")
   }

   var tempIdx = 0
   var workersocre = 0.0
   var score = 0.0

   for Idx, val := range offer.Volunteer {

      grade, _ := strconv.ParseFloat(val.Grade, 64)
      cnt, _ := strconv.ParseFloat(val.Count, 64)
      if grade == 0 || cnt == 0 {
         score = 0.0
      } else {
         score = grade / cnt
      }
      if score >= workersocre {
         tempIdx = Idx
         workersocre = score
      }
   }

   worker := offer.Volunteer[tempIdx]

   var hirdWorkd []Document

   hirdWorkd = append(hirdWorkd, worker)

   offer.Volunteer = hirdWorkd

   offer.State = false

   updateofferBytes, _ := json.Marshal(offer)
   APIstub.PutState(offerkey, updateofferBytes)

   return shim.Success([]byte("Success setFreelancerd"))
}

// 지원자 검증
// params - 이력서 키, 구인공고 키
func (s *SmartContract) verifyConditions(APIstub shim.ChaincodeStubInterface, args []string) pb.Response {
   mspid, _ := cid.GetMSPID(APIstub)

   if mspid != "CertOrg" {
      return shim.Error("You do not have permission")
   }

   if len(args) != 2 {
      return shim.Error("Incorrect number of arguments. Expecting 2")
   }

   docskey, offerkey := args[0], args[1]

   // cc1 invoke
   verifyConditionsResponse := s.returnCertificate(APIstub, []string{docskey})
   if verifyConditionsResponse.GetStatus() >= 400 {
      return shim.Error(fmt.Sprintf("failed to bringDocument err : %s", verifyConditionsResponse.GetMessage()))
   }

   freelancerDoc := verifyConditionsResponse.GetPayload()

   document := Document{}
   json.Unmarshal(freelancerDoc, &document)

   // check Condition

   offerBytes, err := APIstub.GetState(offerkey)
   if err != nil {
      return shim.Error("Could not find Offer")
   }
   offer := Offer{}
   json.Unmarshal(offerBytes, &offer)

   for key, values := range offer.Requirement {
      val, exists := document.Certification[key]
      if key == "TOEIC" {
         offval, _ := strconv.Atoi(values)
         teerval, _ := strconv.Atoi(val)
         if offval > teerval {
            return shim.Error("You lack the grade of " + key)
         }
         continue
      } else if key == "Rating" {
         wokeCount, _ := strconv.ParseFloat(document.Count, 64)
         workgrade, _ := strconv.ParseFloat(document.Grade, 64)

         rate, _ := strconv.ParseFloat(val, 64)

         var grade float64 = workgrade / wokeCount

         if rate > grade {
            return shim.Error("Bye")
         }
         continue
      } else if key == "Count" {
         wantcount, _ := strconv.Atoi(val)
         workcount, _ := strconv.Atoi(document.Count)
         if wantcount > workcount {
            return shim.Error("Bye")
         }
         continue
      } else if key == "OPIC" {
         opicGrade := [7]string{"NL", "NM", "NH", "IL", "IM", "IH", "AL"}
         var workgrade int
         var wantgrade int
         for idx, grade := range opicGrade {
            if grade == values {
               wantgrade = idx
            }
            if grade == val {
               workgrade = idx
            }
         }
         if wantgrade > workgrade {
            return shim.Error("Byte")
         }
         continue
      }
      if !exists {
         return shim.Error("You don't have " + key)
      }
   }

   offer.Volunteer = append(offer.Volunteer, document)

   updateofferBytes, _ := json.Marshal(offer)
   err3 := APIstub.PutState(offerkey, updateofferBytes)
   if err3 != nil {
      return shim.Error("Failed to regist")
   }

   return shim.Success([]byte("Regist Success"))
}

// params - 회사지갑 키, 평점, 구인공고 키, 문서 키
func (s *SmartContract) setRating(APIstub shim.ChaincodeStubInterface, args []string) pb.Response {
   if len(args) != 5 {
      return shim.Error("Incorrect number of arguments. Expecting 5")
   }

   CompanyWallet, psword, workerGrade, offerkey, docsKey := args[0], args[1], args[2], args[3], args[4]

   hash := sha256.New()

   hash.Write([]byte(psword))

   word := hash.Sum(nil)

   psStr := hex.EncodeToString(word)

   check, _ := strconv.Atoi(workerGrade)

   if check > 5 && check < 0 {
      return shim.Error("0<= grade <=5")
   }

   offer := Offer{}
   offerBytes, err := APIstub.GetState(offerkey)
   if err != nil {
      return shim.Error("Failed to Offer getState")
   } else if offerBytes == nil {
      return shim.Error("Could not find Offer")
   }
   json.Unmarshal(offerBytes, &offer)

   if psStr != offer.C_PSWORD {
      return shim.Error("This is not your jobPosting")
   }

   document := Document{}
   docBytes, err := APIstub.GetState(docsKey)
   if err != nil {
      return shim.Error("Failed to Document GetState")
   } else if offerBytes == nil {
      return shim.Error("Could not find document")
   }
   json.Unmarshal(docBytes, &document)

   if document.PSWORD != offer.Volunteer[0].PSWORD {
      return shim.Error("This is not the person you hired")
   }

   // if !offer.State {
   //    return shim.Error("you not hired worker")
   // }

   // 평점 업데이트
   result := s.updateRating(APIstub, []string{docsKey, CompanyWallet, workerGrade, offer.Remuneration})
   if result.GetStatus() >= 400 {
      return shim.Error("Failed to setRating")
   }

   // 구인공고 완전 삭제
   setratingResponse := s.deleteJobPosting(APIstub, []string{offerkey})
   if setratingResponse.GetStatus() >= 400 {
      return shim.Error("Failed to DeleteOffer")
   }

   return shim.Success([]byte("setRating Success"))
}

func (s *SmartContract) getJobPosting(APIstub shim.ChaincodeStubInterface) pb.Response {

   // Find latestKey
   offerkeyAsBytes, _ := APIstub.GetState("offerlatesKey")
   offerkey := OfferKey{}
   json.Unmarshal(offerkeyAsBytes, &offerkey)
   idxStr := strconv.Itoa(offerkey.Idx + 1)

   var startKey = "postKey0"
   var endKey = offerkey.Key + idxStr
   fmt.Println(startKey)
   fmt.Println(endKey)

   resultsIter, err := APIstub.GetStateByRange(startKey, endKey)
   if err != nil {
      return shim.Error(err.Error())
   }

   defer resultsIter.Close()

   var buffer bytes.Buffer
   buffer.WriteString("[")
   bArrayMemberAlreadyWritten := false
   for resultsIter.HasNext() {
      queryResponse, err := resultsIter.Next()
      if err != nil {
         return shim.Error(err.Error())
      }

      // 고용된 공고 건너뜀
      if strings.Contains(string(queryResponse.Value), "false") {
         continue
      }

      if bArrayMemberAlreadyWritten == true {
         buffer.WriteString(",")
      }
      buffer.WriteString("{\"Key\":")
      buffer.WriteString("\"")
      buffer.WriteString(queryResponse.Key)
      buffer.WriteString("\"")

      buffer.WriteString(", \"Record\":")

      buffer.WriteString(string(queryResponse.Value))
      buffer.WriteString("}")
      buffer.WriteString("\n")
      bArrayMemberAlreadyWritten = true
   }
   buffer.WriteString("]\n")
   return shim.Success(buffer.Bytes())
}

// param - 구인공고 키
func (s *SmartContract) deleteJobPosting(APIstub shim.ChaincodeStubInterface, args []string) pb.Response {
   if len(args) != 1 {
      return shim.Error("Incorrect number of arguments. Expecting 1")
   }

   offerKey := args[0]

   // Delete the key from the state in ledger
   err := APIstub.DelState(offerKey)
   if err != nil {
      return shim.Error("Failed to delete state")
   }

   return shim.Success(nil)
}

// param - 이력서 키, 구인공고 키
func (s *SmartContract) setApply(APIstub shim.ChaincodeStubInterface, args []string) pb.Response {
   mspid, _ := cid.GetMSPID(APIstub)

   if mspid != "CertOrg" {
      return shim.Error("You do not have permission")
   }

   if len(args) != 3 {
      return shim.Error("Incorrect number of arguments. Expecting 3")
   }

   workerDockey, psword, offerKey := args[0], args[1], args[2]

   hash := sha256.New()

   hash.Write([]byte(psword))

   word := hash.Sum(nil)

   psStr := hex.EncodeToString(word)

   documentAsBytes, err := APIstub.GetState(workerDockey)
   if err != nil {
      fmt.Println(err.Error())
   } else if documentAsBytes == nil {
      return shim.Error("Could not Find document")
   }

   document := Document{}
   json.Unmarshal(documentAsBytes, &document)

   if psStr != document.PSWORD {
      return shim.Error("You do not have permission")
   }

   result := s.verifyConditions(APIstub, []string{workerDockey, offerKey})
   if result.GetStatus() >= 400 {
      return shim.Error(fmt.Sprintf("failed to registe err : %s", result.GetMessage()))
   }

   return shim.Success([]byte("Volunteer Success"))
}
