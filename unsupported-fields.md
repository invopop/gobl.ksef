# Unsupported fields (hardcoded or not mapped)

Here are features not supported by the converter.

## Hardcoded values

| XML field | Struct field | Value | Notes |
| --------- | ---------- | ----- | ----- |
| `KodFormularza` | `FormCode` | `FA` | |
| `KodFormularza@kodSystemowy` (attribute) | `FormCode.SystemCode` | `FA (3)` | |
| `KodFormularza@wersjaSchemy` (attribute) | `FormCode.SchemaVersion` | `1-0E` | |
| `WariantFormularza` | `FormVariant` | `3` | |
| `SystemInfo` | `SystemInfo` | `Invopop` | |
| `Fa>NoweSrodkiTransportu>P_22N` | `NoNewTransportIntraCommunitySupply` | `1` | For new transport intra-community supply, set `P_22` to 1 (rare special case), otherwise set `P_22N` to 1 |
| `Fa>P_23` | `SimplifiedProcedureBySecondTaxpayer` | `2` | For simplified procedure by second taxpayer (for three-party transactions inside the European Union), set `P_23` to 1, otherwise set `P_23` to 2 |
| `Fa>Platnosc>RachunekBankowyFaktora` | `FactorBankAccounts` | `[]` | Bank account of the factor (third party) |

## Not mapped

The following fields are now present in the structs but are not currently being mapped from GOBL data:

### Seller (Podmiot1)
| XML field | Struct field | Notes |
| --------- | ------------ | ----- |
| `Podmiot1>NrEORI` | `EORI` | EORI number of the seller. The EORI number is the number in the EU Economic Operators Registration and Identification Number |
| `Podmiot1>AdresKoresp` | `CorrespondenceAddress` | Correspondence address if different from main address |
| `Podmiot1>StatusInfoPodatnika` | `TaxpayerStatus` | Taxpayer status: 1=liquidation, 2=restructuring, 3=bankruptcy, 4=inheritance |
| `Adres>GLN` | `GLN` | Global Location Number. GLN is a number that enables, among other things, the identification of physical or functional units within a company. |

### Buyer (Podmiot2)
| XML field | Struct field | Notes |
| --------- | ------------ | ----- |
| `Podmiot2>IDNabywcy` | `BuyerID` | Unique key linking buyer data in correction invoices |
| `Podmiot2>NrEORI` | `EORI` | EORI number of the buyer |
| `Podmiot2>AdresKoresp` | `CorrespondenceAddress` | Correspondence address if different from main address |
| `Podmiot2>NrKlienta` | `CustomerNumber` | Customer number used in contracts/orders |

### Third Party (Podmiot3)
| XML field | Struct field | Notes |
| --------- | ------------ | ----- |
| `Podmiot3>IDNabywcy` | `BuyerID` | Unique buyer link key |
| `Podmiot3>NrEORI` | `EORI` | EORI number |
| `Podmiot3>Rola>RolaInna` | `OtherRole` | Marker for custom role |
| `Podmiot3>Rola>OpisRoli` | `OtherRoleDescription` | Custom role description |
| `Podmiot3>Udzial` | `Share` | Percentage share (e.g., for multiple buyers) |
| `Podmiot3>NrKlienta` | `CustomerNumber` | Customer number |

### Authorized Entity (PodmiotUpowazniony) - COMPLETE STRUCTURE NOT MAPPED
| XML field | Struct field | Notes |
| --------- | ------------ | ----- |
| `PodmiotUpowazniony` | `AuthorizedEntity` | For enforcement authorities, bailiffs, tax representatives |
| `PodmiotUpowazniony>NrEORI` | `EORI` | EORI number |
| `PodmiotUpowazniony>DaneIdentyfikacyjne` | | Identification data (NIP, Name) |
| `PodmiotUpowazniony>Adres` | `Address` | Address |
| `PodmiotUpowazniony>AdresKoresp` | `CorrespondenceAddress` | Correspondence address |
| `PodmiotUpowazniony>DaneKontaktowe>EmailPU` | `Email` | Email |
| `PodmiotUpowazniony>DaneKontaktowe>TelefonPU` | `Phone` | Phone |
| `PodmiotUpowazniony>RolaPU` | `Role` | Role: 1=enforcement authority, 2=bailiff, 3=tax representative |

### Invoice (Fa)
| XML field | Struct field | Notes |
| --------- | ------------ | ----- |
| `Fa>P_6` | `CompletionDate` | The date of delivery or completion of the delivery of goods or services or the date of receipt of payment, referred to in Art. 106b sec. 1(4) of the Act, if such date is specified and differs from the date of issue of the invoice.|
| `Fa>WZ` | `WarehouseDocuments` | Warehouse document numbers (0-1000) |
| `Fa>KursWalutyZ` | `CurrencyRateForTax` | Exchange rate for tax calculation |
| `Fa>P_15ZK` | `AmountBeforeCorrection` | Amount before correction (for KOR_ZAL and other corrections) |
| `Fa>ZaliczkaCzesciowa` | `PartialAdvancePayments` | Partial advance payments data (array, 0-31) for invoices documenting receipt of multiple payments |
| `Fa>ZaliczkaCzesciowa>P_6Z` | `PaymentDate` | Date of receiving payment |
| `Fa>ZaliczkaCzesciowa>P_15Z` | `PaymentAmount` | Payment amount |
| `Fa>ZaliczkaCzesciowa>KursWalutyZW` | `CurrencyExchangeRate` | Currency exchange rate for tax calculation |
| `Fa>FakturaZaliczkowa` | `AdvanceInvoices` | References to preceding advance invoices (for ROZ type) |
| `Fa>TP` | `TP` | Existing relationships between buyer and supplier of goods or services |
| `Fa>ZwrotAkcyzy` | `ExciseTaxRefund` | Excise tax refund marker for farmers |

### Correction Invoice Fields
| XML field | Struct field | Notes |
| --------- | ------------ | ----- |
| `DaneFaKorygowanej>OkresFaKorygowanej` | `CorrectionPeriod` | Period for discount/reduction corrections |
| `DaneFaKorygowanej>NrFaKorygowany` | `CorrectedInvoiceNo` | Correct invoice number (when fixing wrong number) |
| `DaneFaKorygowanej>Podmiot1K` | `CorrectedSeller` | Seller data from corrected invoice (if changed) |
| `DaneFaKorygowanej>Podmiot2K` | `CorrectedBuyer` | Buyer data from corrected invoice (if changed) |

### Settlement (Rozliczenie) - COMPLETE STRUCTURE NOT MAPPED
| XML field | Struct field | Notes |
| --------- | ------------ | ----- |
| `Fa>Rozliczenie` | `Settlement` | Additional charges and deductions |
| `Fa>Rozliczenie>Obciazenia` | `Charges` | Charges added to total (0-100) |
| `Fa>Rozliczenie>Obciazenia>SumaObciazen` | `TotalCharges` | Sum of all charges |
| `Fa>Rozliczenie>Odliczenia` | `Deductions` | Deductions from total (0-100) |
| `Fa>Rozliczenie>Odliczenia>SumaOdliczen` | `TotalDeductions` | Sum of all deductions |
| `Fa>Rozliczenie>DoZaplaty` | `AmountToPay` | Final amount to pay |
| `Fa>Rozliczenie>DoRozliczenia` | `AmountToSettle` | Overpaid amount to settle/refund |

### Transaction Conditions (WarunkiTransakcji) - PARTIALLY MAPPED
| XML field | Struct field | Notes |
| --------- | ------------ | ----- |
| `Fa>WarunkiTransakcji>Umowy` | `Contracts` | Contract references (date & number, 0-100) |
| `Fa>WarunkiTransakcji>NrPartiiTowaru` | `BatchNumbers` | Product batch numbers (0-1000) |
| `Fa>WarunkiTransakcji>WarunkiDostawy` | `DeliveryTerms` | Incoterms delivery conditions |
| `Fa>WarunkiTransakcji>KursUmowny` | `ContractRate` | Contract exchange rate |
| `Fa>WarunkiTransakcji>WalutaUmowna` | `ContractCurrency` | Contract currency |
| `Fa>WarunkiTransakcji>Transport` | `Transport` | Transport/shipping details (0-20) |
| `Fa>WarunkiTransakcji>PodmiotPosredniczacy` | `IntermediaryParty` | Intermediary entity marker |

**Note**: `Fa>WarunkiTransakcji>Zamowienia` (order references) is now mapped bidirectionally via `Ordering.Purchases`. The first order reference date and number are converted. Other WarunkiTransakcji fields remain unmapped.

### Transport - COMPLETE STRUCTURE NOT MAPPED
| XML field | Struct field | Notes |
| --------- | ------------ | ----- |
| `Transport>RodzajTransportu` | `TransportType` | Transport type: 1=sea, 2=rail, 3=road, 4=air, 5=postal, 7=fixed, 8=waterway, 9=own |
| `Transport>TransportInny` | `OtherTransportType` | Marker for other transport type |
| `Transport>OpisInnegoTransportu` | `OtherTransportDesc` | Description of other transport type |
| `Transport>Przewoznik` | `Carrier` | Carrier information |
| `Transport>NrZleceniaTransportu` | `TransportOrderNumber` | Transport order number |
| `Transport>OpisLadunku` | `CargoType` | Cargo/packaging type code |
| `Transport>LadunekInny` | `OtherCargoType` | Marker for other cargo type |
| `Transport>OpisInnegoLadunku` | `OtherCargoDesc` | Description of other cargo type |
| `Transport>JednostkaOpakowania` | `PackagingUnit` | Packaging unit description |
| `Transport>DataGodzRozpTransportu` | `TransportStartTime` | Transport start date/time |
| `Transport>DataGodzZakTransportu` | `TransportEndTime` | Transport end date/time |
| `Transport>WysylkaZ` | `ShipFrom` | Shipping from address |
| `Transport>WysylkaPrzez` | `ShipVia` | Intermediate shipping addresses (0-20) |
| `Transport>WysylkaDo` | `ShipTo` | Shipping to address |

### Annotations - Extended Fields
| XML field | Struct field | Notes |
| --------- | ------------ | ----- |
| `Fa>Adnotacje>NoweSrodkiTransportu` | `NewTransportMeans` | Complete new transport means structure |
| `Fa>Adnotacje>NoweSrodkiTransportu>P_22` | `Marker` | New transport means marker |
| `Fa>Adnotacje>NoweSrodkiTransportu>P_42_5` | `Art42Obligation` | Art. 42 ust. 5 obligation |
| `Fa>Adnotacje>NoweSrodkiTransportu>NowySrodekTransportu` | `NewTransportMeansItems` | Vehicle/watercraft/aircraft details (0-10000) |

### Line Items (FaWiersz) - Extended Fields
| XML field | Struct field | Notes |
| --------- | ------------ | ----- |
| `FaWiersz>P_6A` | `CompletionDate` | Completion date for this specific line |
| `FaWiersz>Indeks` | `InternalCode` | Internal product code (max 50 chars) |
| `FaWiersz>GTIN` | `GTIN` | Global Trade Item Number (max 20 chars) |
| `FaWiersz>PKWiU` | `PKWiU` | Polish Classification of Products and Services |
| `FaWiersz>CN` | `CN` | Combined Nomenclature code |
| `FaWiersz>PKOB` | `PKOB` | Polish Classification of Construction Objects |
| `FaWiersz>P_11A` | `GrossPriceTotal` | Gross total price (for art. 106e ust. 7-8) |
| `FaWiersz>P_11Vat` | `VATAmount` | VAT amount (for art. 106e ust. 10) |
| `FaWiersz>P_12_Zal_15` | `Attachment15GoodsMarker` | Split payment marker (value: 1) |
| `FaWiersz>KursWaluty` | `CurrencyRate` | Currency exchange rate for this line |

### Order Line Items (ZamowienieWiersz) - Extended Fields
| XML field | Struct field | Notes |
| --------- | ------------ | ----- |
| `ZamowienieWiersz>UU_IDZ` | `UniqueID` | Universal unique order line ID |
| `ZamowienieWiersz>IndeksZ` | `InternalCode` | Internal product code |
| `ZamowienieWiersz>GTINZ` | `GTIN` | Global Trade Item Number |
| `ZamowienieWiersz>PKWiUZ` | `PKWiU` | Polish Classification code |
| `ZamowienieWiersz>CNZ` | `CN` | Combined Nomenclature code |
| `ZamowienieWiersz>PKOBZ` | `PKOB` | Construction objects code |

### Payment - Extended Fields
| XML field | Struct field | Notes |
| --------- | ------------ | ----- |
| `Platnosc>ZaplataCzesciowa>PlatnoscInna` | `OtherPaymentMeanMarker` | Marker for other payment method |
| `Platnosc>ZaplataCzesciowa>OpisPlatnosci` | `OtherPaymentMean` | Description of other payment method |
| `Platnosc>TerminPlatnosci>TerminOpis` | `TermDescription` | Alternative textual payment deadline format |
| `Platnosc>LinkDoPlatnosci` | `PaymentLink` | Payment link URL with IPKSeF parameter |
| `Platnosc>IPKSeF` | `KSeFPaymentID` | KSeF payment identifier (13 chars) |

### Other Not Mapped Fields
| XML field | Notes |
| --------- | ----- |
| `Stopka` | Footer information - not required in schema, identifies parties in national databases |
| `Zalacznik` | Attachment structure for custom data (key-value pairs or tables) - not required |

## Unset fields

The following fields are present in the struct to be converted to XML, but are not set anywhere in the GOBL → KSeF direction:

| XML field | Struct field | Notes |
| --------- | ------------ | ----- |
| `Fa>P_1M` | `IssuePlace` | issue place - not required in schema |
| `Fa>P_14_1W` | `StandardRateTaxConvertedToPln` | for foreign currency invoices |
| `Fa>P_14_2W` | `ReducedRateTaxConvertedToPln` | for foreign currency invoices |
| `Fa>P_14_3W` | `SuperReducedRateTaxConvertedToPln` | for foreign currency invoices |
| `Fa>FP` | `FP` | indicates a case where an invoice is issued in addition to a regular receipt - not required in schema |
| `Fa>FaWiersz>StanPrzed` | `BeforeCorrectionMarker` | in a correction invoice, indicates that the line describes the state before the correction. Used during KSeF → GOBL parsing but not set during GOBL → KSeF building. |
| `Fa>FaWiersz>GTU` | `SpecialGoodsCode` | Code identifying certain classes of goods and services (01 = alcoholic beverages, 02 = vehicle fuels...), values GTU_01 to GTU_13 |

## Other cases

- `P_12` now handles both numeric rates (e.g. `23`, `8`) and text values (`0 WDT`, `0 EX`, `zw`, `np I`, `np II`, `oo`, `0 KR`) for zero-rated, exempt, and special tax categories.

## How the listing is done

The following check is done by this way:
- Take [official sample XML files provided by the Polish authorities](https://www.gov.pl/attachment/937002fa-c6b5-477d-8b56-22105fa728c2)
- For each XML field in the reference file, find place in our code that generates it
- If a reference file contains a field and we don't have a corresponding field in our code, we mark it as "Not mapped".
- If a field is present in our struct, but there exist one possible value for this field, such a field goes to "Hardcoded values" list.
- If a field is present in our struct, but there exist no code that sets this field, such a field goes to "Unset values" list.

Note that this particular check:
- Does not cover all possible fields in the XML schema, only the ones present in the sample files
- Does not detect fields in our code that are not present in the sample files, or in the XML schema at all

### List of test cases from the sample pack

1. FP invoice (invoice added to receipt)
2. Correction invoice with `StanPrzed` - two lines, one with `StanPrzed` set to 1 (shows the state before the correction), another without (after correction)
3. Correction invoice - with a single line containing negative numbers (alternative method to case 2)
4. Regular invoice with a third party, with multiple delivery dates
5. Correction invoice where customer's data changed, contains `Podmiot2K` and `IDNabywcy` fields, and no line items (`FaWiersz`) - it's only for correcting customer's identification data and nothing else
6. Correction invoice to give a discount, for multiple earlier corrected invoices (`DaneFaKorygowanej`), contains `OkresFaKorygowanej` field, and no line items (`FaWiersz`)
7. Differs from example 6 by having a line item (`FaWiersz`) - meaning that the invoice corrects only part of the earlier deliveries, not all of them
8. Sale of used goods - uses margin scheme, contains partial payment, contains `P_PMarzy_3_1` field, contains `P_11A` field
9. Invoice to a local government unit, contains both tax exempt items (`P_13_7`) and standard tax items (`P_13_1`)
10. Advance invoice, already paid in full, with two customers, each having 50% share
11. Corrects invoice from example 10 because of incorrect used tax rate, type `KOR_ZAL`, contains negative numbers for the incorrect tax rate in fields `P_13_1` and `P_14_1`, contains `P_13_2` and `P_14_2` fields for the correct tax rate
12. Similar to example 11, but corrects amount of advance payment.
13. Similar to example 11, but corrects total payment amount.
14. Settlement invoice (`ROZ`) finalizing a series of advance invoices. Finalizes transaction from examples 10, 11, 12, 13.
15. Simplified invoice (`UPR`) containing `P_13_1` and `P_14_1` fields (total net amount and total tax amount)
16. Simplified invoice (`UPR`) not containing these fields, but total net amount and tax amount is possible to calculate from the line items (field `P_12` containing tax percentage) - in this case, this is allowed, but the converter does not need to handle this alternative
17. Similar to example 14 but containing incorrectly calculated values
18. Correction to a settlement invoice (type `KOR_ROZ`) - corrects invoice from example 17
19. Invoice using margin scheme, contains `P_PMarzy` and `P_PMarzy_2` fields, meaning that the seller is a travel agency
20. Regular invoice, but using EUR instead of PLN, contains `KursWaluty` (currency exchange rate) for each line item, and `P_14_1W` for totals converted to PLN. Also contains `WarunkiTransakcji` (transaction conditions) with transport-related information inside.
21. Similar to example 20, but each line item has different currency exchange rate (even though the whole invoice is in EUR).
22. Intra-EU supply, tax-exempt - contains `P_13_6_2` field, customer is identified by EU VAT number (`KodUE` + `NrVatUE`)
23. Export outside European Union, tax-exempt - contains `P_13_6_3` field, customer is identified by respective country's tax identifier (`KodKraju` + `NrID`)
24. Invoice containing an attachment, missing values in a table are marked with `-`
25. Invoice containing an attachment, missing values in a table are marked with empty XML element (alternative to example 24)
26. Invoice containing additional charges (in the example, additional charges are for package recycling) - contains `Rozliczenie` field

## References

What is FP:
- http://www.vademecumpodatnika.pl/artykul_narzedziowa,1190,0,19945,ujmowanie-w-ewidencji-vat-faktur-do-paragonow.html

GTU codes:
- https://www.podatki.biz/artykuly/jpkvat-z-deklaracja-oznaczenia-dostawy-i-swiadczenia-uslug-gtu_4_45232.htm

When should `P_11A` be used - art. 106e ust. 7 i 8:
- https://sip.lex.pl/akty-prawne/dzu-dziennik-ustaw/podatek-od-towarow-i-uslug-17086198/art-106-e

When should `P_23` be used:
- https://isp-modzelewski.pl/serwis/wewnatrzwspolnotowe-transakcje-trojstronne/#:~:text=ramach%20procedury%20uproszczonej.-,Zgodnie%20z%20art.,dodanej%20ostatniego%20w%20kolejno%C5%9Bci%20podatnika.
- https://www.krgroup.pl/procedura-uproszczona-rozliczenia-vat-w-wewnatrzwspolnotowej-transakcji-trojstronnej/#:~:text=Warunki%20zastosowania%20procedury,realizowanej%20w%20ramach%20procedury%20uproszczonej.
- https://sip.lex.pl/akty-prawne/dzu-dziennik-ustaw/podatek-od-towarow-i-uslug-17086198/dz-12-roz-8

When should `P_15ZK` be used (art. 106f ust. 3) for remaining amount to pay before correction:
https://sip.lex.pl/akty-prawne/dzu-dziennik-ustaw/podatek-od-towarow-i-uslug-17086198/art-106-f
